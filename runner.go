package gomtr

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gogather/com"
	"github.com/gogather/safemap"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const maxttls = 50

// service
type MtrService struct {
	taskQueue *safemap.SafeMap
	flag      int64
	index     int64
	in        io.WriteCloser
	out       io.ReadCloser
	outChan   chan string
}

func NewMtrService() *MtrService {
	return &MtrService{
		taskQueue: safemap.New(),
		flag:      102400,
		index:     1,
		in:        nil,
		out:       nil,
		outChan:   make(chan string, 1000),
	}
}

// start service and wait mtr-packet stdio
func (ms *MtrService) Start() {
	go ms.startup()
	time.Sleep(time.Second)
}

func (ms *MtrService) startup() {

	cmd := exec.Command("./mtr-packet")

	var e error

	ms.out, e = cmd.StdoutPipe()
	if e != nil {
		fmt.Println(e)
	}

	ms.in, e = cmd.StdinPipe()
	if e != nil {
		fmt.Println(e)
	}

	err, e := cmd.StderrPipe()
	if e != nil {
		fmt.Println(e)
	}

	// start sub process
	if e := cmd.Start(); nil != e {
		fmt.Printf("ERROR: %v\n", e)
	}

	// read data and put into result chan
	go func() {
		for {
			// read lines
			bio := bufio.NewReader(ms.out)
			for{
				output, isPrefix, err := bio.ReadLine()
				if err != nil {
					break
				}

				if string(output) != "" {
					ms.outChan <- string(output)
				}

				if isPrefix {
					break
				}
			}


		}
	}()

	// get result from chan and parse
	go func() {
		for {
			select {
			case result := <-ms.outChan:
				{
					ms.parseTTLData(result)
				}
			}

		}
	}()

	// error output
	go func() {
		for {
			var readBytes []byte = make([]byte, 100)
			err.Read(readBytes)
			time.Sleep(1)
		}
	}()

	// wait sub process
	if e := cmd.Wait(); nil != e {
		fmt.Printf("ERROR: %v\n", e)
	}

}

func (ms *MtrService) send(id int64, ip string, c int) {
	defer func() {
		recover()
	}()

	if c > 100 {
		c = 99
	} else if c < 1 {
		c = 1
	}

	for i := 1; i <= c; i++ {
		sendId := id*10000 + int64(i)*100
		for idx := 1; idx <= maxttls; idx++ {
			ms.in.Write([]byte(fmt.Sprintf("%d send-probe ip-4 %s ttl %d\n", sendId+int64(idx), ip, idx)))
		}
	}

}

func (ms *MtrService) Request(ip string, c int, callback func(interface{})) {

	task := &MtrTask{
		id:       ms.index,
		callback: callback,
		c:        c,
		ttlData:  safemap.New(),
	}

	ms.taskQueue.Put(fmt.Sprintf("%d", ms.index), task)

	ms.send(ms.index, ip, c)

	ms.index++

	if ms.index > ms.flag {
		ms.index = 1
	}

}

func (ms *MtrService) parseTTLData(data string) {
	segs := strings.Split(data, "\n")

	for i := 0; i < len(segs); i++ {
		item := strings.TrimSpace(segs[i])
		if len(item) > 0 {
			ms.parseTTLDatum(item)
		}
	}
}

func (ms *MtrService) parseTTLDatum(data string) {
	// what i got
	fmt.Println(data)

	hasNewline := strings.Contains(data, "\n")
	if hasNewline {
		fmt.Println(hasNewline)
	}

	segments := strings.Split(data, " ")

	var ttlData *TTLData
	var fullID int64
	var ttlTime int64
	//var err error

	if len(segments) >= 1 {
		idInt, err := strconv.Atoi(segments[0])
		if err != nil {
			idInt = 0
		}
		fullID = int64(idInt)
	}

	if len(segments) >= 2 {
		if segments[1] == "command-parse-error" {
			ttlData = &TTLData{
				TTLID: ms.getTTLID(fullID),
				err:   errors.New("command parse error"),
			}
		} else if segments[1] == "no-reply" {
			ttlData = &TTLData{
				TTLID: ms.getTTLID(fullID),
				err:   errors.New("no reply"),
			}
		}
	}

	if len(segments) >= 6 {
		ttlTimeInt, err := strconv.Atoi(segments[5])
		if err != nil {
			ttlTimeInt = 0
		} else {
			ttlTime = int64(ttlTimeInt)
		}

		ttlData = &TTLData{
			TTLID:  ms.getTTLID(fullID),
			ipType: segments[2],
			ip:     segments[3],
			time:   ttlTime,
		}
	}

	// store
	taskID := fmt.Sprintf("%d", ms.getRealID(fullID))
	task, ok := ms.taskQueue.Get(taskID)
	if ok {
		ttlID := ms.getTTLID(fullID)
		task.(*MtrTask).save(ttlID, ttlData)
	}

	// check task
	if ok {
		if task.(*MtrTask).check() {
			// callback
			cb := task.(*MtrTask).callback
			if cb != nil {
				cb(task.(*MtrTask))
				task.(*MtrTask).clear()
				ms.taskQueue.Remove(taskID)
			}
		}
	}

}

func (ms *MtrService) getTTLID(fullID int64) int {
	idStr := fmt.Sprintf("%d", fullID)
	length := len(idStr)
	ttlStr := com.SubString(idStr, length-4, 4)
	ttl, e := strconv.Atoi(ttlStr)
	if e != nil {
		ttl = 0
	}
	return ttl
}

func (ms *MtrService) getRealID(fullID int64) int64 {
	return fullID / 10000
}
