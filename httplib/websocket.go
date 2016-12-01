package httplib

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

//类型定义
//当连接成功时
type OnConnectCall func()

//当接收到消息时
type OnMessageCall func()

//当关闭连接时
type OnCloseCall func()

//当遇到错误时
type OnErrorCall func()
type WebSocket struct {
}

func (obj WebSocket) New() *WebSocket {
	return &obj
}

type Wsocket struct {
	lock          *sync.Mutex
	Conns         map[string]*websocket.Conn
	onConnect     func(sbid string)
	onMessageCall func(sbid string, msg []byte)
	onCloseCall   func(sbid string)
	grader        websocket.Upgrader
}

func NewWebSocket() *Wsocket {
	o := new(Wsocket)
	o.grader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//CheckOrigin:     func(r *http.Request) bool { return true },
	}
	o.Conns = make(map[string]*websocket.Conn) //初始化连接池
	o.lock = new(sync.Mutex)
	return o
}

//sbid=识别id
func (this *Wsocket) Hanlder(sbid string, w http.ResponseWriter, r *http.Request) (err error) {
	this.lock.Lock()
	if conn, ok := this.Conns[sbid]; ok {
		conn.Close()
	}
	this.Conns[sbid], err = this.grader.Upgrade(w, r, nil)
	this.lock.Unlock()
	if err != nil {
		return
	}
	this.onConnect(sbid)
	for {
		_, msg, err := this.Conns[sbid].ReadMessage()
		if err != nil {
			this.onCloseCall(sbid)   //
			delete(this.Conns, sbid) //删除连接池
			break
		}
		this.onMessageCall(sbid, msg)
	}
	return
}

//当第一次连接时
func (this *Wsocket) OnConnect(call func(sbid string)) {
	this.onConnect = call
}

//当消息上来时
func (this *Wsocket) OnMessage(call func(sbid string, msg []byte)) {
	this.onMessageCall = call
}

//当消息断开时
func (this *Wsocket) OnClose(call func(sbid string)) {
	this.onCloseCall = call
}

//发送消息
func (this *Wsocket) Write(sbid string, msg string) error {
	if conn, ok := this.Conns[sbid]; ok {
		if err := conn.WriteMessage(1, []byte(msg)); err != nil {
			return errors.New(sbid + "已掉线")
		}
	}
	return nil
}

//关闭连接
func (this *Wsocket) Close(sbid string) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	e := this.Conns[sbid].Close()
	delete(this.Conns, sbid)
	return e
}

//给所有在线的发送消息
func (this *Wsocket) SendMessageAll(msg string) (err error) {
	msgStr := []byte(msg)
	for _, conn := range this.Conns {
		err = conn.WriteMessage(1, msgStr)
		if err != nil {
			continue
		}
	}
	return nil
}

//计算多少人在线
func (this *Wsocket) CountConnects() int {
	return len(this.Conns)
}
