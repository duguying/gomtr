package gomtr

import (
	"testing"
	"fmt"
)

func Test_Mtr(t *testing.T) {
	mtr := NewMtrService()
	mtr.Start()

	mtr.Request("duguying.net",2, func() {
		fmt.Println("hello, mtr")
	})
}
