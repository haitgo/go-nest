package websocket

import (
	"go-nest/cache"

	"github.com/gorilla/websocket"
)

//websocket的连接对象
type Conn struct {
	conn         *websocket.Conn //
	ConnId       string          //连接编号
	cache.Memory                 //继承缓存包
}

//实例化一个连接对象
func (obj Conn) New() *Conn {
	return &obj
}

//发送消息
func (this *Conn) Write(data []byte) (err error) {
	return this.conn.WriteMessage(1, data)
}

//关闭连接
func (this *Conn) Close() (err error) {
	return this.conn.Close()
}
