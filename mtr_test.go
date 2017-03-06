package gomtr

import (
	"fmt"
	"testing"
	"time"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService()
	go mtr.Start()

	mtr.Request("183.131.7.130", 2, func() {
		fmt.Println("hello, mtr")
	})

	for {
		time.Sleep(1)
	}
}
