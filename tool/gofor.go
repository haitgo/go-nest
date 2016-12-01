package tool

import (
	"time"
)

//go协程管理器，用于死循环的关闭
type Gofor struct {
	state bool
}

//tm=时间(time.Duration 或者 string)
func NewGofor(tm interface{}, call func()) *Gofor {
	var t time.Duration
	if ttm, ok := tm.(time.Duration); ok {
		t = ttm
	} else if ttm, ok := tm.(string); ok {
		t, _ = time.ParseDuration(ttm)
	}
	if t <= 0 {
		return nil
	}
	obj := new(Gofor)
	obj.state = true
	go func() {
		for {
			if !obj.state {
				return
			}
			time.Sleep(t)
			call()
		}
	}()
	return obj
}

//关闭这个死循环
func (this *Gofor) Stop() {
	this.state = false
}
