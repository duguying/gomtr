package gomtr

import (
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

// service
type MtrService struct {
	taskQueue *safemap.SafeMap
	flag      int64
	index     int64
	in        io.WriteCloser
	out       io.ReadCloser
}

func NewMtrService() *MtrService {
	return &MtrService{
		taskQueue: safemap.New(),
		flag:      10240000,
		index:     1,
		in:        nil,
		out:       nil,
	}
}

func (ms *MtrService) Start() {

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

	// read data
	go func() {
		for {
			var readBytes []byte = make([]byte, 100)
			ms.out.Read(readBytes)
			fmt.Print(string(readBytes))
			time.Sleep(1)
		}
	}()

	go func() {
		for {
			// todo: use channel for output
			var readBytes []byte = make([]byte, 100)
			err.Read(readBytes)
			ms.parseTTLData(string(readBytes))
			time.Sleep(1)
		}
	}()

	if e := cmd.Start(); nil != e {
		fmt.Printf("ERROR: %v\n", e)
	}

	if e := cmd.Wait(); nil != e {
		fmt.Printf("ERROR: %v\n", e)
	}

}

func (ms *MtrService) send(id int64, ip string, ttls int) {
	sendId := id * 100
	for idx := 1; idx <= ttls; idx++ {
		ms.in.Write([]byte(fmt.Sprintf("%d send-probe ip-4 %s ttl %d\r", sendId+int64(idx), ip, idx)))
	}
}

func (ms *MtrService) Request(ip string, ttls int, callback func()) {

	task := mtrTask{
		id:       ms.index,
		callback: callback,
		ttls:     ttls,
		ttlData:  safemap.New(),
	}

	ms.taskQueue.Put(fmt.Sprintf("%d", ms.index), task)

	ms.send(ms.index, ip, ttls)

	ms.index++

}

func (ms *MtrService) parseTTLData(data string) {
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
	task, ok := ms.taskQueue.Get(fmt.Sprintf("%d", ms.getRealID(fullID)))
	if ok {
		ttlID := ms.getTTLID(fullID)
		task.(*mtrTask).save(ttlID, ttlData)
	}

	// check task
	if ok {
		if task.(*mtrTask).check() {
			// callback
			cb := task.(*mtrTask).callback
			if cb != nil {
				cb()
			}
		}
	}

}

func (ms *MtrService) getTTLID(fullID int64) int {
	idStr := fmt.Sprintf("%d", fullID)
	length := len(idStr)
	ttlStr := com.SubString(idStr, length-2, 2)
	ttl, e := strconv.Atoi(ttlStr)
	if e != nil {
		ttl = 0
	}
	return ttl
}

func (ms *MtrService) getRealID(fullID int64) int64 {
	return fullID / 100
}
