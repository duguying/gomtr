package gomtr

import (
	"fmt"
	"github.com/gogather/safemap"
	"strconv"
)

// parsed ttl item data
type TTLData struct {
	TTLID  int
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

func (mt *MtrTask) check() bool {
	for idx := 1; idx <= mt.c; idx++ {
		for i := 1; i <= maxttls; i++ {
			idstr := fmt.Sprintf("%02d%02d", idx, i)
			id, _ := strconv.Atoi(idstr)
			_, ok := mt.ttlData.Get(fmt.Sprintf("%d", id))
			if !ok {
				return false
			}
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
			if ok {
				ttlid := itemData.TTLID
				cid := ttlid / 100
				ttl := ttlid % 100
				_, ok := results[ttl]
				if !ok {
					results[ttl] = map[int]int64{}
				}
				results[ttl][cid] = itemData.time
			}
		}
	}
	return results
}
