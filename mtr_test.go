package gomtr

import (
	"testing"
	"time"
	"github.com/gogather/com/log"
	"github.com/gogather/com"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService()
	go mtr.Start()

	i := 1
	for {
		mtr.Request("183.131.7.130",10, func(response interface{}) {
			//fmt.Println(response)
			task:=response.(*MtrTask)
			log.Blueln(com.JsonEncode(task.GetResult()))
		})
		i++
		time.Sleep(time.Second)
	}

	for {
		time.Sleep(1)
	}
}
