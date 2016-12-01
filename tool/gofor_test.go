package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestGofor(t *testing.T) {
	d := NewGofor("1s", func() {
		fmt.Println("a")
	})
	time.Sleep(time.Second * 2)
	d.Stop()
	time.Sleep(time.Second * 10)

}
