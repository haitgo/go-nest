package websocket

import (
	"go-nest/cache"
	"go-nest/tool"
	"net/http"

	"github.com/gorilla/websocket"
)

//websocket封装库
type WebSocket struct {
	cache.Memory                //连接池
	onConnectCall OnConnectCall //当连接时触发
	onMessageCall OnMessageCall //当消息上来时触发
	onCloseCall   OnCloseCall   //当关闭连接是触发
	onErrorCall   OnErrorCall   //当发生错误是触发
	grader        websocket.Upgrader
}

func (obj WebSocket) New() *WebSocket {
	obj.grader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		//CheckOrigin:     func(r *http.Request) bool { return true },
	}
	return &obj
}

//注册当websocket连接时回调
func (this *WebSocket) OnConnect(call OnConnectCall) {
	this.onConnectCall = call
}

//注册当websocket接收到消息时的回调
func (this *WebSocket) OnMessage(call OnMessageCall) {
	this.onMessageCall = call
}

//注册当websocket关闭连接时的回调
func (this *WebSocket) OnClose(call OnCloseCall) {
	this.onCloseCall = call
}

//注册当websocket发生错误时的回调
func (this *WebSocket) OnError(call OnErrorCall) {
	this.onErrorCall = call
}

//处理webSocket请求
//输入：connId=连接编号
func (this *WebSocket) Handle(connId string, w http.ResponseWriter, r *http.Request) (err error) {
	//读取到一个连接
	conn, err := this.grader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	objConn := Conn{ConnId: connId, conn: conn}.New()
	this.Set(connId, objConn) //保存连接到连接池
	this.onConnectCall(objConn)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			this.onCloseCall(objConn)
			break
		}
		tool.Try(func() {
			this.onMessageCall(objConn, msg)
		}, func(e error) {
			this.onErrorCall(objConn, e)
		})
	}
	return nil
}

//给指定的连接发送消息
func (this *WebSocket) WriteTo(connectId interface{}, data []byte) (err error) {
	var conn *Conn
	this.GetBy(connectId, conn)
	return conn.Write(data)
}

//关闭指定连接
func (this *WebSocket) CloseTo(connectId interface{}) (err error) {
	var conn *Conn
	this.GetBy(connectId, conn)
	return conn.Close()
}

//给所有的连接发送消息
func (this *WebSocket) WriteToAll(data []byte) {
	this.Each(func(k, v interface{}) {
		if conn, ok := v.(*Conn); ok {
			conn.Write(data)
		}
	})
}
