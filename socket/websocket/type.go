package websocket

//类型定义
//当连接成功时
type OnConnectCall func(c *Conn)

//当接收到消息时
type OnMessageCall func(c *Conn, data []byte)

//当关闭连接时
type OnCloseCall func(c *Conn)

//当遇到错误时
type OnErrorCall func(c *Conn, e error)
