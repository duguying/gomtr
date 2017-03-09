package gomtr

import (
	"fmt"
	"github.com/gogather/safemap"
	"io"
	"sort"
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
	raw    string
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
	start := time.Now().UnixNano() / 1000000

	for {
		// get tllID
		tllID := getTTLID(rid)

		// check ready
		d, ok := mt.ttlData.Get(fmt.Sprintf("%d", tllID))
		if !ok || d == nil {
			// not ready, continue
		} else {
			data, ok := d.(*TTLData)
			if !ok || data == nil {
				// not ready, continue
			} else {
				// ready, check replied
				if data.status == "ttl-expired" || data.err != nil {
					// not get replied
					return 1
				} else if data.status == "reply" {
					// get replied
					return 0
				}
			}
		}

		now := time.Now().UnixNano() / 1000000

		// timeout
		if now > start {
			return 1
		}

		time.Sleep(1)
	}

	// this will not reached
	return 1
}

func (mt *MtrTask) clear() {
	for key, _ := range mt.ttlData.GetMap() {
		mt.ttlData.Remove(key)
	}
}

func (mt *MtrTask) GetResultMap() map[int]map[int]int64 {
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

func (mt *MtrTask) GetSummary() map[int]map[string]string {
	results := map[int][]*TTLData{}

	var keys []int
	for ks := range mt.ttlData.GetMap() {
		k, _ := strconv.Atoi(ks)
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, key := range keys {
		item, ok := mt.ttlData.Get(fmt.Sprintf("%d", key))
		if ok {
			itemData, ok := item.(*TTLData)
			if ok && itemData != nil {
				ttlid := itemData.TTLID
				ttl := ttlid % 100

				// put data
				array, ok := results[ttl]
				if ok {
					results[ttl] = append(array, itemData)
				} else {
					results[ttl] = []*TTLData{itemData}
				}
			}
		}
	}

	summarys := map[int]map[string]string{}
	for key, value := range results {
		// summary
		summarys[key] = map[string]string{
			"Last":  fmtNumber(sortLast(value)),
			"Avg":   fmtNumber(sortAvg(value)),
			"Best":  fmtNumber(sortBest(value)),
			"Wrst":  fmtNumber(sortWorst(value)),
			"StDev": fmtNumber(sortSTDev(value)),
			"Snt":   fmt.Sprintf("%d", sortSnt(value)),
			"IP":    sortLastTTLData(value).ip,
			"ttl":   fmt.Sprintf("%d", key),
		}
	}
	return summarys
}

func fmtNumber(n float64) string {
	return fmt.Sprintf("%1.01f", n/1000)
}

func sortLastTTLData(array []*TTLData) *TTLData {
	l := len(array)
	return array[l-1]
}

func sortLast(array []*TTLData) float64 {
	for i := len(array) - 1; i >= 0; i-- {
		item := array[i]
		if item.err == nil {
			return float64(item.time)
		}
	}
	return 0
}

func sortSnt(array []*TTLData) int {
	if len(array) <= 0 || array == nil {
		return 0
	}

	return len(array)
}

func sortAvg(array []*TTLData) float64 {
	if len(array) <= 0 || array == nil {
		return 0
	}

	var result int64 = 0
	c := 0
	for i := 0; i < len(array); i++ {
		item := array[i]
		if item.err == nil {
			result = result + item.time
			c++
		}
	}
	return float64(result) / float64(c)
}

func sortBest(array []*TTLData) float64 {
	var best int64 = -1

	for i := 0; i < len(array); i++ {
		item := array[i]
		if item.err == nil {
			if item.time > best {
				best = item.time
			}
		}
	}

	return float64(best)
}

func sortWorst(array []*TTLData) float64 {
	var worst int64 = array[0].time

	for i := 0; i < len(array); i++ {
		item := array[i]
		if item.err == nil {
			if item.time < worst {
				worst = item.time
			}
		}
	}

	return float64(worst)
}

// 标准偏差统计
func sortSTDev(array []*TTLData) float64 {
	if len(array) <= 0 || array == nil {
		return 0
	}

	avg := sortAvg(array)

	var s float64 = 0
	var c int = 0
	for i := 0; i < len(array); i++ {
		item := array[i]
		if item.err == nil {
			c++
			itemTime := float64(item.time)
			d := itemTime - avg
			delta := d * d
			s = s + delta
		}
	}

	return s / float64(c-1)
}
