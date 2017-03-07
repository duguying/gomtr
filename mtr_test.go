package gomtr

import (
	"fmt"
	"testing"
	"time"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService()
	go mtr.Start()

	i := 1
	for {
		mtr.Request("183.131.7.130",10, func(response interface{}) {
			fmt.Println(response)
		})
		i++
	}

	for {
		time.Sleep(1)
	}
}
