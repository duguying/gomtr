package gomtr

import (
	"bufio"
	"errors"
	"fmt"
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
	taskQueue     *safemap.SafeMap
	flag          int64
	index         int64
	in            io.WriteCloser
	out           io.ReadCloser
	outChan       chan string
	mtrPacketPath string
}

// NewMtrService new a mtr service
// path - mtr-packet executable path
func NewMtrService(path string) *MtrService {
	return &MtrService{
		taskQueue:     safemap.New(),
		flag:          102400,
		index:         1,
		in:            nil,
		out:           nil,
		outChan:       make(chan string, 1000),
		mtrPacketPath: path,
	}
}

// Start start service and wait mtr-packet stdio
func (ms *MtrService) Start() {
	go ms.startup()
	time.Sleep(time.Second)
}

func (ms *MtrService) startup() {

	cmd := exec.Command(ms.mtrPacketPath)

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
			for {
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

// Request send a task request
// ip       - the test ip
// c        - repeat time, such as mtr tool argument c
// callback - just callback after task ready
func (ms *MtrService) Request(ip string, c int, callback func(interface{})) {

	if c <= 0 {
		c = 1
	}

	task := &MtrTask{
		id:       ms.index,
		callback: callback,
		c:        c,
		ttlData:  safemap.New(),
	}

	ms.index++

	ms.taskQueue.Put(fmt.Sprintf("%d", ms.index), task)

	task.send(ms.in, ms.index, ip, c)

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

	hasNewline := strings.Contains(data, "\n")
	if hasNewline {
		fmt.Println(hasNewline)
	}

	segments := strings.Split(data, " ")

	var ttlData *TTLData
	var fullID int64
	var ttlTime int64
	var ttlerr error
	var status string
	var ipType string
	var ip string

	if len(segments) <= 0 {
		return
	}

	if len(segments) > 0 {
		idInt, err := strconv.Atoi(segments[0])
		if err != nil {
			idInt = 0
		}
		fullID = int64(idInt)
	}

	if len(segments) > 1 {
		switch segments[1] {
		case "command-parse-error", "no-reply", "probes-exhausted", "network-down", "permission-denied", "no-route", "invalid-argument", "feature-support":
			{
				ttlerr = errors.New(segments[1])
				break
			}
		case "ttl-expired":
			{
				status = segments[1]
				break
			}
		case "reply":
			{
				status = segments[1]
				break
			}
		}

	}

	if len(segments) > 2 {
		ipType = segments[2]
	}

	if len(segments) > 3 {
		ip = segments[3]
	}

	if len(segments) > 5 {
		ttlTimeInt, err := strconv.Atoi(segments[5])
		if err != nil {
			ttlTimeInt = 0
		} else {
			ttlTime = int64(ttlTimeInt)
		}
	}

	ttlData = &TTLData{
		TTLID:        getTTLID(fullID),
		ipType:       ipType,
		ip:           ip,
		err:          ttlerr,
		status:       status,
		raw:          data,
		time:         ttlTime,
		receivedTime: time.Now(),
	}

	// store
	taskID := fmt.Sprintf("%d", getRealID(fullID))
	taskRaw, ok := ms.taskQueue.Get(taskID)
	var task *MtrTask = nil
	if ok && taskRaw != nil {
		task = taskRaw.(*MtrTask)
		ttlID := getTTLID(fullID)
		task.save(ttlID, ttlData)
	} else {
		return
	}

}
