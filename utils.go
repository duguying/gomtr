package gomtr

import (
	"fmt"
	"github.com/gogather/com"
	"math"
	"sort"
	"strconv"
	"time"
)

func fmtNumber(n float64) string {
	return fmt.Sprintf("%.01f", n/1000)
}

func sortHasReply(array []*TTLData) bool {
	for i := len(array) - 1; i >= 0; i-- {
		item := array[i]
		if item.status == "reply" {
			return true
		}
	}
	return false
}

// least array all is no-reply
func sortLeastAllNoReply(least []*TTLData) bool {
	for i := len(least) - 1; i >= 0; i-- {
		item := least[i]
		if item.status != "no-reply" {
			return false
		}
	}
	return true
}

func clearSummary(summary map[int]map[string]string) {
	var ttlKeys []int
	for k := range summary {
		ttlKeys = append(ttlKeys, k)
	}
	sort.Ints(ttlKeys)

	for i := len(ttlKeys) - 1; i >= 0; i-- {
		j := i - 1
		if j >= 0 {
			if summary[i]["ip"] == "???" && summary[j]["ip"] == "???" {
				delete(summary, i)
			}
		}
	}
}

func sortLastTTLData(array []*TTLData) *TTLData {
	l := len(array)
	for i := l - 1; i >= 0; i-- {
		item := array[i]
		if item.err == nil {
			return item
		}
	}

	return array[l-1]
}

func sortIP(array []*TTLData) string {
	last := sortLastTTLData(array)
	if last.err != nil {
		return "???"
	} else {
		return last.ip
	}
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

func sortSntReality(array []*TTLData) int {
	if len(array) <= 0 || array == nil {
		return 0
	}

	c := 0
	for i := 0; i < len(array); i++ {
		item := array[i]
		if item.err == nil {
			c++
		}
	}

	return c
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

	if c <= 0 {
		c = 1
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

	if best < 0 {
		best = 0
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

	if c <= 0 {
		c = 1
	}

	stdev := math.Sqrt(s / float64(c-1))
	if math.IsNaN(stdev) {
		stdev = float64(0)
	}

	return stdev
}

func getTTLID(fullID int64) int {
	idStr := fmt.Sprintf("%d", fullID)
	length := len(idStr)
	ttlStr := com.SubString(idStr, length-4, 4)
	ttl, e := strconv.Atoi(ttlStr)
	if e != nil {
		ttl = 0
	}
	return ttl
}

func getRealID(fullID int64) int64 {
	return fullID / 10000
}

func getMtrStartTime() string {
	now := time.Now()
	week := com.SubString(now.Weekday().String(), 0, 3)
	month := com.SubString(now.Month().String(), 0, 3)
	return fmt.Sprintf("%s %s %d %s", week, month, now.Day(), now.Format("15:04:05 2006"))
}
