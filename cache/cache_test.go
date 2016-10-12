package cache

import (
	"testing"
)

func Test_Memory(t *testing.T) {
	c := New(Memory(300))
	c.Set("a", "王该帖")
	t.Log(c.GetString("a"))
	//time.Sleep(time.Second * 3)
	c.Expire(100)
	t.Log("3秒后：", c.GetString("a"))
}
func Test_Redis(t *testing.T) {
	c := New(Redis("localhost:6379", "", 1, "user"))
	//
	c.Set("a", "王该帖")
	t.Log(c.GetString("a"))
	//time.Sleep(time.Second * 3)
	c.Expire(100)
	t.Log("3秒后：", c.GetString("a"), c.Len())
}
