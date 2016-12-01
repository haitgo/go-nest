package cache

import (
	"encoding/json"
	"errors"
	rds "go-nest/Redis"
	"strconv"
	"sync"
)

var (
	clients     = make(map[string]*rds.Client) //Redis客户端池
	clientsLock = new(sync.Mutex)              //map互斥锁，保证数据读取安全
)

//此处使用的是Redis的hash，如果想用Redis的特殊功能，请使用Redis包
type Redis struct {
	client  *rds.Client
	keyName string //该名称用于存储Redis的hash的key名
	expire  int64  //延迟时间
}

//Redis缓存
//addr数据库地址，db使用的数据库，pswd密码
func NewRedis(addr, pswd string, db int, keyName string) *Redis {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	obj := new(Redis)
	k := addr + pswd + strconv.Itoa(db)
	if _, ok := clients[k]; !ok {
		clients[k] = &rds.Client{
			Addr:     addr,
			Db:       db,
			Password: pswd,
		}
	}
	obj.client = clients[k]
	obj.keyName = keyName
	return obj
}

//延时,所有的项目延时
func (this *Redis) Expire(exp int64) {
	this.expire = exp
	keys, err := this.client.Keys(this.keyName + "*")
	if err != nil {
		return
	}
	for _, key := range keys {
		this.client.Expire(key, exp)
	}
}

//添加值
func (this *Redis) Set(key, value interface{}, exp ...int64) (err error) {
	kata, err := json.Marshal(key) //key序列化
	if err != nil {
		return errors.New("键名输入错误")
	}
	newkey := this.keyName + string(kata)
	data, err := json.Marshal(value) //value序列化
	if err != nil {
		return err
	}
	err = this.client.Set(newkey, data)
	if err != nil {
		return err
	}
	expire := this.expire
	if len(exp) == 1 {
		expire = exp[0]
	}
	if expire != 0 {
		this.client.Expire(newkey, expire)
	}
	this.client.Hset(this.keyName, newkey, []byte{1})
	return nil
}

//生成Redis的key
func (this *Redis) compileKey(key interface{}) string {
	kata, err := json.Marshal(key) //key序列化
	if err != nil {
		return ""
	}
	return this.keyName + string(kata)
}

//返回值：正确[]byte{}，错误nil
func (this *Redis) Get(key interface{}) interface{} {
	newkey := this.compileKey(key) //key序列化
	data, err := this.client.Get(newkey)
	if err != nil {
		return nil
	}
	var vl interface{}
	err = json.Unmarshal(data, &vl)
	if err != nil {
		return nil
	}
	return vl
}

//获取值为int类型
func (this *Redis) GetInt(key interface{}) int {
	if i, ok := this.Get(key).(int); ok {
		return i
	}
	return 0
}

//获取值为float类型
func (this *Redis) GetFloat(key interface{}) float64 {
	if i, ok := this.Get(key).(float64); ok {
		return i
	}
	return 0
}

//获取值为string类型
func (this *Redis) GetString(key interface{}) string {
	if i, ok := this.Get(key).(string); ok {
		return i
	}
	return ""
}

//获取值，并自动转换赋值给value
func (this *Redis) GetBy(key, value interface{}) {
	newkey := this.compileKey(key) //key序列化
	data, err := this.client.Get(newkey)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, value)
	if err != nil {
		return
	}
}

//删除
func (this *Redis) Del(key interface{}) {
	newkey := this.compileKey(key) //key序列化
	if newkey != "" {
		this.client.Del(newkey)
		this.client.Hdel(this.keyName, newkey)
	}
}

//清空
func (this *Redis) Clean() {
	keys, err := this.client.Keys(this.keyName + "*")
	if err != nil {
		return
	}
	for _, newkey := range keys {
		this.client.Del(newkey)
		this.client.Hdel(this.keyName, newkey)
	}
}

//保存在Redis的缓存数据长度，如果超时则删除该缓存
func (this *Redis) Len() int {
	length, err := this.client.Hlen(this.keyName)
	if err != nil {
		return 0
	}
	return length
}

//遍历缓存
func (this *Redis) Each(call func(key, value interface{})) {
	keys, err := this.client.Keys(this.keyName + "*")
	if err != nil {
		return
	}
	for _, key := range keys {
		data, err := this.client.Get(key)
		if err != nil {
			continue
		}
		kata, err := this.client.Hget(this.keyName, key)
		if err != nil {
			continue
		}
		call(kata, data)
	}
}
