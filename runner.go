package main

import (
	"os/exec"
	"time"
	"fmt"
)

func main() {
	cmd := exec.Command("./mtr-packet")

	out,e := cmd.StdoutPipe()
	if e != nil {
		fmt.Println(e)
	}

	in,e:=cmd.StdinPipe()
	if e != nil {
		fmt.Println(e)
	}

	err,e:=cmd.StderrPipe()
	if e != nil {
		fmt.Println(e)
	}

	go func() {
		for  {
			var readBytes []byte = make([]byte,100)
			out.Read(readBytes)
			fmt.Print(string(readBytes))
			time.Sleep(time.Second)
		}

	}()

	go func() {
		for{
			in.Write([]byte("1 send-probe ip-4 183.131.7.130 ttl 1\n"))
			time.Sleep(time.Second)
		}

	}()

	go func() {
		for  {
			var readBytes []byte
			err.Read(readBytes)
			fmt.Print(string(readBytes))
			time.Sleep(time.Second)
		}
	}()


	if e := cmd.Start(); nil != e {
		fmt.Printf("ERROR-: %v\n", e)
	}
	if e := cmd.Wait(); nil != e {
		fmt.Printf("ERROR+: %v\n", e)
	}
}
