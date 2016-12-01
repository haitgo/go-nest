package bind

import (
	"testing"
)

type User struct {
	Name  int    `valid:"String" `
	Phone string `valid:"Phone"`
}

func TestBind(t *testing.T) {
	data := map[string]string{
		"Name":  "132",
		"Phone": "155282032s20",
	}
	obj := new(User)
	err, key := Map(data).Need("Name", "Phone").Set(obj)
	if err != nil {
		t.Log(key, "值绑定错误")
	}
	t.Log("绑定结果是：", obj)
}
