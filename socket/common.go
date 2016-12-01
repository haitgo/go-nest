package socket

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	TYPE_REQUEST  int8 = 1
	TYPE_RESPONSE int8 = 2
)

var (
	DEBUG        = false
	MSG_HEAD     = []byte{19, 88}            //包头
	MSG_HEAD_STR = string(MSG_HEAD)          //发送的报头
	WRITE_WAIT   = time.Nanosecond * 5000000 //写入等待时间
)

//设置为调试模式,会打印调试信息
func SetDebugModel() {
	DEBUG = true
}

//调试输出
func debug(arg ...interface{}) {
	if DEBUG {
		now := time.Now().Format("[2006-01-02 15:04:05]")
		def := []interface{}{now}
		fmt.Println(append(def, arg...)...)
	}
}

//服务器端当连接时回调函数，必须返回一个数据作为连接池编号
//输入：data=client发起连接时发送的初次连接数据
type ServerConnectHandleFunc func(c *Context) interface{}

//客户端连接成功时回调
//输入：data=服务器端连接成功后返回的数据
type ClientConnectHandleFunc func(c *Context)

//当关闭连接时的回调
type CloseHandleFunc func(c *Context)

//当发生错误时的回调
type ErrorHandleFunc func(c *Context)

//响应回调
//返回处理的结果
type ResponseHandleFunc func(c *Context) interface{}

//请求回调
type RequestHandleFunc func(c *Context)

//socket事件回调
//通讯数据
type MsgData struct {
	Type      int8   //请求类型 1=request,2=response
	Route     string //路由
	Timestemp int64  //请求时间戳，保证回调调用时只使用一次
	Data      []byte //数据内容
}

//消息回调的上下文
type Context struct {
	Conn  *Conn  //连接对象
	Data  []byte //数据
	Error error
}

func newContext(c *Conn, data []byte, e error) *Context {
	o := new(Context)
	o.Conn = c
	o.Data = data
	o.Error = e
	return o
}

//将数据转换格式
func (this *Context) Unmarshal(v interface{}) error {
	return json.Unmarshal(this.Data, v)
}

//
func (this *Context) Set(k, v interface{}, exp ...int64) error {
	return this.Conn.Set(k, v, exp...)
}

//
func (this *Context) Get(k interface{}) interface{} {
	return this.Conn.Get(k)
}

//
func (this *Context) Del(k interface{}) {
	this.Conn.Del(k)
}

func (this *Context) Write(dt interface{}) (int, error) {
	return this.Conn.Write(dt)
}

//发送byte类型数据
func (this *Context) WriteByte(data []byte) (n int, err error) {
	return this.Conn.conn.Write(data)
}

//
func (this *Context) Close() error {
	return this.Conn.Close()
}
