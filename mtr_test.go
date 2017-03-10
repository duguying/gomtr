package gomtr

import (
	"fmt"
	"github.com/gogather/com/log"
	"testing"
	"time"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService("./mtr-packet")
	go mtr.Start()

	time.Sleep(time.Second * 10)

	iplist := []string{"183.131.7.130", "127.0.0.1", "114.215.151.25", "111.13.101.208"}

	for i := 0; i < len(iplist); i++ {
		mtr.Request(iplist[i], 10, func(response interface{}) {
			task := response.(*MtrTask)
			log.Blueln("[ID]", task.id)
			fmt.Println(task.GetSummaryString())
		})
	}

	for {
		time.Sleep(1)
	}
}
