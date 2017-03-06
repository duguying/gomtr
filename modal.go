package gomtr

import (
	"github.com/gogather/safemap"
	"fmt"
)

// parsed ttl item data
type TTLData struct {
	TTLID  int
	ipType string
	ip     string
	time   int64
	err    error
}

// task
type mtrTask struct {
	id       int64
	callback func()
	ttls     int
	ttlData  *safemap.SafeMap // item is ttlData, key is ttl
}

func (mt *mtrTask) save(ttl int, data *TTLData) {
	mt.ttlData.Put(fmt.Sprintf("%d", ttl), data)
}

func (mt *mtrTask) check() bool {
	for i := 1; i <= mt.ttls; i++ {
		_, ok := mt.ttlData.Get(fmt.Sprintf("%d", i))
		if !ok {
			return false
		}
	}
	return true
}