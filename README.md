# gomtr [![Build Status](https://travis-ci.org/duguying/gomtr.svg?branch=master)](https://travis-ci.org/duguying/gomtr)

gomtr is a golang wrap for mtr-packet with born for solve concurrency mtr calling.

### usage

```golang
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

	time.Sleep(time.Second * 10)

	iplist := []string{"183.131.7.130", "127.0.0.1", "114.215.151.25", "111.13.101.208"}

	for i := 0; i < len(iplist); i++ {
		mtr.Request(iplist[i], 10, func(response interface{}) {
			task := response.(*gomtr.MtrTask)
			log.Blueln("[ID]", i)
			fmt.Println(task.GetSummaryString())
		})
	}
}

```

### license

MIT License
