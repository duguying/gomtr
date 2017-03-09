package gomtr

import (
	"fmt"
	"github.com/gogather/com"
	"math"
	"strconv"
)

func fmtNumber(n float64) string {
	return fmt.Sprintf("%.01f", n/1000)
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

func sortSntReality(array []*TTLData) int {
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
