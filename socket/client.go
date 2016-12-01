package socket

import (
	"encoding/json"
	"go-nest/tool"
	"net"
)

//socket客户端
type Client struct {
	Conn                                //继承连接对象
	onConnect   ClientConnectHandleFunc //回调函数
	onClose     CloseHandleFunc         //回调函数
	onError     ErrorHandleFunc         //回调函数
	connectData []byte                  //连接数据
}

//创建一个客户端
//输入：data=初次连接输入的参数
func NewClient(data interface{}) *Client {
	dt, ok := data.([]byte)
	if !ok {
		dt, _ = json.Marshal(data)
	}
	obj := new(Client)
	obj.connectData = dt
	obj.Conn.init()
	return obj
}

//注册回调函数，当发起连接成功时触发
func (this *Client) OnConnect(call ClientConnectHandleFunc) {
	this.onConnect = call
}

//注册回调函数，当接收到消息时触发
func (this *Client) OnMessage(route string, call ResponseHandleFunc) {
	if route != "" {
		this.response.Set(route, call)
	}
	this.Conn.defResponseCall = call

}

//注册回调函数，当关闭连接时触发
func (this *Client) OnClose(call CloseHandleFunc) {
	this.onClose = call
}

//注册回调函数，当发生错误时触发
func (this *Client) OnError(call ErrorHandleFunc) {
	this.onError = call
}

//发起连接
//输入：addr=发起请求的地址
func (this *Client) Connect(addr string) error {
	c, err := net.Dial("tcp", addr)
	if err != nil { //连接失败
		return err
	}
	this.Conn.conn = c
	defer func() {
		this.Conn.Close()
	}()
	tool.Try(func() {
		this.Write(this.connectData)
		connData, err := this.Read()
		if err != nil {
			this.onError(newContext(&this.Conn, nil, err))
			return
		}
		this.onConnect(newContext(&this.Conn, connData, nil))
		for {
			data, err := this.Read()
			if err != nil {
				this.onClose(newContext(&this.Conn, data, err))
				return
			}
			go tool.Try(func() { //协程执行每个请求
				this.HandleRouter(data)
			}, func(e error) {
				this.onError(newContext(&this.Conn, nil, e))
			})
		}
	}, func(err error) {
		this.onError(newContext(&this.Conn, nil, err))
		this.Close()
	})
	return nil
}
