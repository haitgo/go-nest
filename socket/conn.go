package socket

import (
	"encoding/json"
	"errors"
	"go-nest/cache"
	"go-nest/tool"
	"net"
	"sync"
	"time"
)

//-----------------------------------------------------------------
//socket连接
type Conn struct {
	Id              interface{}        //当前连接的编号
	LastTime        time.Time          //最后访问时间
	conn            net.Conn           //net默认的连接对象
	defResponseCall ResponseHandleFunc //默认响应
	request         *cache.Memory      //发出请求
	response        *cache.Memory      //接收请求
	reqIndex        int64              //请求的序号
	lostBuf         []byte             //未读完的缓冲数据
	sync.Mutex                         //继承读写锁
	cache.Memory                       //继承内存缓存作为保存常用数据
}

func NewConn(c net.Conn) *Conn {
	o := new(Conn)
	o.conn = c //保存连接对象
	o.init()
	return o
}
func (this *Conn) init() {
	this.request = cache.NewMemory(1)
	this.response = cache.NewMemory(1)
}

//读取通讯数据
//输入：rl=读取初始数据长度
func (this *Conn) Read() ([]byte, error) {
	var buf []byte
	buffer := make([]byte, 102400)
	for {
		if len(this.lostBuf) > 12 {
			debug("receive stick package.")
			buf = this.parseRead(this.lostBuf, len(this.lostBuf))
		} else {
			this.lostBuf = make([]byte, 0)
			l, err := this.conn.Read(buffer)
			if err != nil {
				return nil, err
			}
			debug("收到数据:", string(buffer[:l]))
			this.LastTime = time.Now() //更新链接时间
			buf = this.parseRead(buffer, l)
		}
		if len(buf) > 0 {
			return buf, nil
		}
	}
	return buf, nil
}

//解析读取的数据
func (this *Conn) parseRead(buffer []byte, l int) []byte {
	newBuf := buffer[:l]
	bufLen := len(newBuf)
	if bufLen <= 5 {
		return newBuf
	}
	head := string(buffer[:2])
	if head != MSG_HEAD_STR { //判断头是否正确
		return newBuf
	}
	length := uint16(0)
	tool.ByteToNumber(buffer[2:4], &length)
	if bufLen < int(length) { //判断数据长度是否足够
		return newBuf
	}
	foot := buffer[length+4]
	if this.m256(4+int(length)) != foot { //判断尾巴是否正确
		return newBuf
	}
	end := int(length + 4)
	if l >= end { //数据是否足够
		buf := buffer[4:end]
		this.lostBuf = buffer[end+1 : l] //剩下的数据
		return buf
	}
	return newBuf
}

//发送数据，dt=任意数据类型
//数据包裹格式：[]byte{a,b,c,d},a=2位报头，b=2位长度，c=数据,d=1位报尾
func (this *Conn) Write(dt interface{}) (n int, err error) {
	time.Sleep(time.Nanosecond * 5000000)
	var data []byte
	if bt, ok := dt.([]byte); ok {
		data = bt
	} else if data, err = json.Marshal(dt); err != nil {
		return 0, err
	}
	btLen := tool.NumberToByte(uint16(len(data)))
	wtdt := append(MSG_HEAD, btLen...)
	wtdt = append(wtdt, data...)
	wtdt = append(wtdt, this.m256(len(wtdt)))
	return this.conn.Write(wtdt)
}

//关闭连接
func (this *Conn) Close() error {
	if this.conn == nil {
		return errors.New("not connected.")
	}
	return this.conn.Close()
}

//
func (this *Conn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *Conn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *Conn) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}

func (this *Conn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}

func (this *Conn) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}

//注册路由
func (this *Conn) RegRouter(route string, call ResponseHandleFunc) {
	this.response.Set(route, call)
}

//默认路由
func (this *Conn) DefRouter(call ResponseHandleFunc) {
	this.defResponseCall = call
}

//调用远端路由,将请求的回调保存到请求池里面
//输入：route=请求的路由，args=请求的参数，call=请求成功后的回调函数
func (this *Conn) CallRouter(route string, arg interface{}, call RequestHandleFunc) {
	this.Lock()
	defer this.Unlock()
	this.reqIndex++
	btArg, _ := json.Marshal(arg)
	data := new(MsgData)
	data.Route = route
	data.Type = TYPE_RESPONSE
	data.Timestemp = this.reqIndex
	data.Data = btArg
	req, dataByte := this.encode(data)
	this.request.Set(string(req), call)
	this.Write(dataByte) //发送请求
	debug("resquest send:", len(dataByte))
}

//阻塞式请求回调，直到有数据返回，可设置超时时间（后期完善）
func (this *Conn) CallRouterSync(route string, arg interface{}, call RequestHandleFunc) {
	this.CallRouter(route, arg, func(c *Context) {
		call(c)
	})
	return
}

//处理请求
func (this *Conn) HandleRouter(d []byte) {
	debug("receive:", len(d))
	this.Lock()
	defer this.Unlock()
	data, err := this.decode(d)
	if err == nil {
		switch data.Type {
		case TYPE_REQUEST: //请求回调
			if call := this.request.Get(data.Route); call != nil {
				call.(RequestHandleFunc)(newContext(this, data.Data, nil)) //
				this.request.Del(data.Route)                               //执行完就删除该事件
				return
			}
		case TYPE_RESPONSE: //处理请求
			if call := this.response.Get(data.Route); call != nil {
				ret := call.(ResponseHandleFunc)(newContext(this, data.Data, nil)) //使用回调处理的结果
				if ret != nil {
					data.Type = TYPE_REQUEST //反馈状态
					var ok bool
					if data.Data, ok = ret.([]byte); !ok {
						data.Data, _ = json.Marshal(ret) //执行结果编码处理
					}
					_, btdt := this.encode(data) //数据编码
					this.Write(btdt)             //发送回调
					debug("response send:", len(btdt))
				}
				return
			}
		}
	}
	//执行默认处理的回调
	if this.defResponseCall != nil {
		this.defResponseCall(newContext(this, d, nil))
	}
	return
}

//发起请求时数据加密
//输出：req=请求+时间戳
//数据结构：{a,b,c,d,e}  a=请求类型(1位），b=时间戳（8位），c=路由长度(8位)，d=路由,e=参数
//响应时，路由=路由+时间
func (this *Conn) encode(d *MsgData) (req []byte, data []byte) {
	btType := tool.NumberToByte(d.Type)
	btTimstemp := tool.NumberToByte(d.Timestemp)
	btRoute := []byte(d.Route)
	btRouteLen := tool.NumberToByte(int64(len(btRoute)))
	data = append(data, btType...)
	data = append(data, btTimstemp...)
	data = append(data, btRouteLen...)
	data = append(data, btRoute...)
	data = append(data, d.Data...)
	//	data = append(data, []byte("\n")...)
	//
	req = append(req, btTimstemp...)
	req = append(req, btRoute...)
	return
}

//收到请求是数据解密
//
func (this *Conn) decode(bt []byte) (data *MsgData, err error) {
	if len(bt) <= 9 {
		return nil, errors.New("data nil.")
	}
	data = new(MsgData)
	tool.ByteToNumber(bt[:1], &data.Type)
	if data.Type != TYPE_REQUEST && data.Type != TYPE_RESPONSE {
		return nil, errors.New("type error.")
	}
	tool.ByteToNumber(bt[1:9], &data.Timestemp)
	if data.Timestemp <= 0 {
		return nil, errors.New("token id error.")
	}
	var rtLen int64
	tool.ByteToNumber(bt[9:17], &rtLen)
	if rtLen <= 0 {
		//print("id is ", rtLen, "\n")
		return nil, errors.New("data error.")
	}
	//print("id is ", rtLen, "\n")
	if data.Type == TYPE_REQUEST {
		rt := bt[1:9]
		rt = append(rt, bt[17:17+rtLen]...)
		data.Route = string(rt)
	}
	if data.Type == TYPE_RESPONSE {
		data.Route = string(bt[17 : 17+rtLen])
	}
	data.Data = bt[17+rtLen:]
	return data, nil
}

func (this *Conn) m256(l int) byte {
	return byte(l % 256)
}
