package socket

//处理socket服务器端

import (
	"go-nest/cache"
	"go-nest/tool"
	"net"
	"time"
)

//tcp服务器
type Server struct {
	onConnect        ServerConnectHandleFunc     //当连接的时候触发的回调
	onClose          CloseHandleFunc             //当关闭连接时触发的回调
	onMessageDefault ResponseHandleFunc          //模式消息处理
	onError          ErrorHandleFunc             //当发生错误时触发的回调
	routerEvent      map[interface{}]interface{} //路由处理
	pool             *cache.Memory               //继承内存缓存作为连接池
}

//创建一个server,并初始化连接池
func NewServer(size int) *Server {
	obj := new(Server)
	obj.pool = cache.NewMemory(size)                    //设置memery的初始内存大小，同时也作为连接池的数量控制
	obj.routerEvent = make(map[interface{}]interface{}) //
	return obj
}

//注册当socket连接时回调,返回一个连接编号
func (this *Server) OnConnect(call ServerConnectHandleFunc) {
	this.onConnect = call
}

//注册当socket接收到消息时的回调
func (this *Server) OnMessage(route string, call ResponseHandleFunc) {
	if route != "" {
		this.routerEvent[route] = call
	}
	this.onMessageDefault = call
}

//注册当socket关闭连接时的回调
func (this *Server) OnClose(call CloseHandleFunc) {
	this.onClose = call
}

//注册当socket发生错误时的回调
func (this *Server) OnError(call ErrorHandleFunc) {
	this.onError = call
}

//通过连接编号获取某个连接对象
func (this *Server) GetConn(id string) (conn *Conn) {
	if c := this.pool.Get(id); c != nil {
		conn = c.(*Conn)
	}
	return
}

//遍历所有连接对象
func (this *Server) EachConn(call func(c *Conn)) {
	this.pool.Each(func(id, conn interface{}) {
		if c, ok := conn.(*Conn); ok {
			call(c)
		}
	})
}

//检查是否超过tm，如果超过则执行call回调
func (this *Server) CheckTimeout(tm int64, call func(c *Conn)) {
	now := time.Now().Unix()
	this.pool.Each(func(id, conn interface{}) {
		if c, ok := conn.(*Conn); ok && now-c.LastTime.Unix() > tm {
			call(c)
		}
	})
}

//发起监听
func (this *Server) Listen(laddr string) error {
	listener, err := net.Listen("tcp", laddr)
	if err != nil {
		return err
	}
	for {
		//等待客户端接入
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go tool.Try(func() {
			this.handleConnect(conn)
		}, func(e error) {
			this.onError(newContext(nil, nil, e))
		})
	}
	return nil
}

//处理每个连接
func (this *Server) handleConnect(c net.Conn) {
	conn := NewConn(c) //创建连接
	defer func() {
		this.onClose(newContext(conn, nil, nil))
		conn.Close()           //保证正常关闭连接
		this.pool.Del(conn.Id) //删除连接
	}()
	conn.response.LoadData(this.routerEvent) //将消息事件写入到每个处理器里面
	conn.DefRouter(this.onMessageDefault)    //添加默认消息处理回调函数
	firstData, err := conn.Read()
	if err != nil {
		this.onError(newContext(conn, nil, err))
		return
	}
	conn.Id = time.Now().Unix() + int64(this.pool.Len())    //使用连接池长度作为地址编号
	ret := this.onConnect(newContext(conn, firstData, nil)) //初次连接发送的数据
	conn.Write(ret)                                         //连接成功后马上返回连接处理后的数据
	this.pool.Set(conn.Id, conn)                            //保存连接到连接池
	for {
		data, err := conn.Read()
		if err != nil {
			break
		}
		go tool.Try(func() {
			conn.HandleRouter(data)
		}, func(e error) {
			this.onError(newContext(conn, data, e))
		})
	}
}

//检查错误
func (this *Server) checkError(err error) int {
	if err != nil {
		if err.Error() == "EOF" {
			return 0
		}
		return -1
	}
	return 1
}
