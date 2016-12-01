package main

import (
	"fmt"
	. "go-nest/socket"
)

type TestStruct struct {
	Name string
	Age  int
}

func main() {
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
