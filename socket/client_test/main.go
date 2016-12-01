package main

import (
	"fmt"
	. "go-nest/socket"
	"strconv"
	"strings"
	"sync"
	"time"
)

var Request = 0
var Response = 0
var Success = 0

type TestStruct struct {
	Name string
	Age  int
}

var lock = new(sync.Mutex)

func main() {
	SetDebugModel()
	cli := NewClient([]byte("292929"))
	cli.OnConnect(func(c *Conn, data []byte) {
		fmt.Println("客户：", string(data))
		go func() {
			str := makeString(10000)
			drx, _ := time.ParseDuration(strconv.Itoa(int(len(str))) + "ns")
			fmt.Println(drx)
			for i := 0; i < 1000; i++ {
				go func() {
					lock.Lock()
					defer lock.Unlock()
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
				}()
				//time.Sleep(time.Nanosecond * 5000000)
			}
			time.Sleep(time.Second * 20)
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
	d := []string{}
	for i := 0; i < l; i++ {
		d = append(d, "王")
	}
	return strings.Join(d, "")
}
