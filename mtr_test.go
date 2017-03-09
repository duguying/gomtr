package gomtr

import (
	"fmt"
	"github.com/gogather/com/log"
	"testing"
	"time"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService()
	go mtr.Start()

	time.Sleep(time.Second * 10)

	i := 1
	for {
		mtr.Request("183.131.7.130", 10, func(response interface{}) {
			//fmt.Println(response)
			task := response.(*MtrTask)
			log.Blueln("[ID]", task.id)
			fmt.Println(task.GetSummaryString())
		})
		i++
		//time.Sleep(time.Second)
		//break
	}

	for {
		time.Sleep(1)
	}
}
