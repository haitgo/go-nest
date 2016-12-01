package cache

import (
	"testing"
	"time"
)

func Test_Memory(t *testing.T) {
	return
	go func() {
		c := NewMemory(300)
		c.Set("a", "王该帖")
		t.Log(c.GetString("a"))
		//time.Sleep(time.Second * 3)
		c.Expire(100)
		c.Each(func(k, v interface{}) {
			c.Del("a")
		})
		//t.Log("3秒后：", c.GetString("a"))
	}()
	time.Sleep(time.Second * 100)
}
func aTest_Redis(t *testing.T) {
	c := NewRedis("localhost:6379", "", 1, "user")
	//
	c.Set("a", "王该帖")
	t.Log(c.GetString("a"))
	//time.Sleep(time.Second * 3)
	c.Expire(100)
	t.Log("3秒后：", c.GetString("a"), c.Len())
}
