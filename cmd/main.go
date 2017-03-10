package main

import (
	"fmt"
	"github.com/duguying/gomtr"
	"github.com/gogather/com/log"
	"time"
)

func main() {
	mtr := gomtr.NewMtrService("./mtr-packet")
	go mtr.Start()

	time.Sleep(time.Second * 5)

	iplist := []string{"183.131.7.130", "127.0.0.1", "114.215.151.25", "111.13.101.208"}

	for i := 0; i < len(iplist); i++ {
		mtr.Request(iplist[i], 10, func(response interface{}) {
			task := response.(*gomtr.MtrTask)
			log.Bluef("[ID] %d cost: %d ms\n", i, task.CostTime / 1000000)
			fmt.Println(task.GetSummaryString())
		})
	}
}
