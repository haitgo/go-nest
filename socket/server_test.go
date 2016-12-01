package socket

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

var Request = 0
var Response = 0
var Success = 0

type TestStruct struct {
	Name string
	Age  int
}

func TestSocket(t *testing.T) {
	go _Server(t)
	time.Sleep(time.Second * 1)
	go _Client(t)
	time.Sleep(time.Second * 2)

}

func _Server(t *testing.T) {
	sev := NewServer(100)
	sev.OnConnect(func(conn *Conn, data []byte) interface{} {
		fmt.Println(string(data), "连接成功")
		return string(data)
	})
	sev.OnMessage("/name", func(conn *Conn, data []byte) interface{} {
		return data
	})
	sev.OnClose(func(conn *Conn) {
		fmt.Println("服务：", "已关闭", conn.Id)
	})
	sev.OnError(func(conn *Conn, e error) {
		if conn != nil {
			fmt.Println("服务：", "遇到意外", e)
		}
		fmt.Println("服务：", "连接失败")
	})
	sev.Listen("localhost:9099")
}

var lock = new(sync.Mutex)

func _Client(t *testing.T) {
	cli := NewClient([]byte("292929"))
	cli.OnConnect(func(c *Conn, data []byte) {
		fmt.Println("客户：", string(data))
		go func() {
			lock.Lock()
			defer lock.Unlock()
			str := makeString(10000)
			drx, _ := time.ParseDuration(strconv.Itoa(int(len(str))) + "ns")
			fmt.Println(drx)
			for i := 0; i < 100; i++ {
				//go func() {

				//fmt.Println("1客户：", "发起测试请求")
				Request++
				d := &TestStruct{Age: 22, Name: str}
				cli.CallRouter("/name", d, func(c *Conn, data []byte) {
					//fmt.Println("5客户：", "返回的结果是", string(data))
					Response++
					t := new(TestStruct)
					c.Unmarshal(data, &t)
					//fmt.Println(t)
					//fmt.Printf("%t \t, %t\n", d, "王海涛")
					if t.Age == 22 {
						Success++
					}
				})
				//}()
				//time.Sleep(time.Nanosecond * 5000000)

			}
			fmt.Println("请求", Request, "次")
			fmt.Println("响应", Response, "次")
			fmt.Println("成功", Success, "次")
		}()

	})
	cli.OnClose(func(c *Conn) {
		fmt.Println("客户：", "连接已经关闭")
	})
	cli.OnError(func(c *Conn, err error) {
		fmt.Println("客户：", "连接遇到错误", err)
	})
	if err := cli.Connect("localhost:9099"); err != nil {
		fmt.Println("客户：", "连接失败", err)
	}
}

func makeString(l int) string {
	r := ""
	for i := 0; i < l; i++ {
		r = r + "王"
	}
	return r
}
