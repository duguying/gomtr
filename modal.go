package gomtr

import (
	"fmt"
	"github.com/gogather/safemap"
	"io"
	"strconv"
	"time"
)

// parsed ttl item data
type TTLData struct {
	TTLID  int
	status string
	ipType string
	ip     string
	time   int64
	err    error
}

func (td *TTLData) String() string {
	return fmt.Sprintf("")
}

// task
type MtrTask struct {
	id       int64
	callback func(interface{})
	c        int
	ttlData  *safemap.SafeMap // item is ttlData, key is ttl
}

func (mt *MtrTask) save(ttl int, data *TTLData) {
	mt.ttlData.Put(fmt.Sprintf("%d", ttl), data)
}

func (mt *MtrTask) send(in io.WriteCloser, id int64, ip string, c int) {
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
		var prevRid int64 = 0
		for idx := 1; idx <= maxttls; idx++ {
			// get reality id
			rid := sendId + int64(idx)

			// the 1st one does not check
			if idx > 1 {
				// sync check status, it will block until ready, and return status
				// ready:
				//       0  get replied
				//       1  not get replied, such as ttl-expired, continue loop
				if mt.checkLoop(prevRid) == 0 {
					break
				}
			}

			prevRid = rid

			in.Write([]byte(fmt.Sprintf("%d send-probe ip-4 %s ttl %d\n", rid, ip, idx)))

			time.Sleep(time.Millisecond)
		}
	}

	// callback
	mt.callback(mt)
	mt.clear()
	// todo: remove this task from task queue

}

// check latest ttl is replied
//    0  [    returned]  ready and get replied
//    1  [    returned]  ready but not replied, go on
// [-1]  [not returned]  not ready, should block
func (mt *MtrTask) checkLoop(rid int64) int {
	for {
		// get tllID
		tllID := getTTLID(rid)

		// check ready
		d, ok := mt.ttlData.Get(fmt.Sprintf("%d", tllID))
		if !ok || d == nil {
			// not ready, continue
			fmt.Println("not ready d")
		} else {
			data, ok := d.(*TTLData)
			if !ok || data == nil {
				// not ready, continue
				fmt.Println("not ready data")
			} else {
				fmt.Println("[ready]", data.status, data.err)
				// ready, check replied
				if data.status == "ttl-expired" {
					// not get replied
					return 1
				} else if data.status == "reply" {
					// get replied
					return 0
				}
			}
		}

		time.Sleep(time.Microsecond)
	}

	// this will not reached
	return 1
}

func (mt *MtrTask) checkCallback() bool {
	for idx := 1; idx <= mt.c; idx++ {
		for i := 1; i <= maxttls; i++ {
			idstr := fmt.Sprintf("%02d%02d", idx, i)
			id, e := strconv.Atoi(idstr)
			if e != nil {
				return false
			}
			d, ok := mt.ttlData.Get(fmt.Sprintf("%d", id))
			if !ok || d == nil {
				return false
			}
			data, ok := d.(*TTLData)
			if !ok || data == nil {
				return false
			}
			//fmt.Printf("[data] %v\n",data)
		}
	}

	return true
}

func (mt *MtrTask) clear() {
	for key, _ := range mt.ttlData.GetMap() {
		mt.ttlData.Remove(key)
	}
}

func (mt *MtrTask) GetResult() map[int]map[int]int64 {
	results := map[int]map[int]int64{}
	for key, _ := range mt.ttlData.GetMap() {
		item, ok := mt.ttlData.Get(key)
		if ok {
			itemData, ok := item.(*TTLData)
			if ok && itemData != nil {
				ttlid := itemData.TTLID
				cid := ttlid / 100
				ttl := ttlid % 100
				_, ok := results[ttl]
				if !ok {
					results[ttl] = map[int]int64{}
				}
				if itemData.err != nil {
					results[ttl][cid] = -1
				} else {
					results[ttl][cid] = itemData.time
				}
			}
		}
	}
	return results
}
