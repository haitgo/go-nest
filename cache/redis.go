package cache

import (
	"encoding/json"
	"errors"
	rds "go-nest/redis"
	"strconv"
	"sync"
)

var (
	clients     = make(map[string]*rds.Client) //redis客户端池
	clientsLock = new(sync.Mutex)              //map互斥锁，保证数据读取安全
)

//此处使用的是redis的hash，如果想用redis的特殊功能，请使用redis包
type redis struct {
	client  *rds.Client
	keyName string //该名称用于存储redis的hash的key名
	expire  int64  //延迟时间
}

//redis缓存
//addr数据库地址，db使用的数据库，pswd密码
func Redis(addr, pswd string, db int, keyName string) *redis {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	o := new(redis)
	k := addr + pswd + strconv.Itoa(db)
	if _, ok := clients[k]; !ok {
		clients[k] = &rds.Client{
			Addr:     addr,
			Db:       db,
			Password: pswd,
		}
	}
	o.client = clients[k]
	o.keyName = keyName
	return o
}

//延时,所有的项目延时
func (this *redis) Expire(exp int64) {
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
func (this *redis) Set(key, value interface{}, exp ...int64) (err error) {
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
	this.client.Hset(this.keyName, newkey, kata)
	return nil
}

//生成redis的key
func (this *redis) compileKey(key interface{}) string {
	kata, err := json.Marshal(key) //key序列化
	if err != nil {
		return ""
	}
	return this.keyName + string(kata)
}

//返回值：正确[]byte{}，错误nil
func (this *redis) Get(key interface{}) interface{} {
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

//获取值，并自动转换赋值给value
func (this *redis) GetBy(key, value interface{}) {
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
func (this *redis) Del(key interface{}) {
	newkey := this.compileKey(key) //key序列化
	if newkey != "" {
		this.client.Del(newkey)
		this.client.Hdel(this.keyName, newkey)
	}
}

//保存在redis的缓存数据长度，如果超时则删除该缓存
func (this *redis) Len() int {
	keys, err := this.client.Keys(this.keyName + "*")
	if err != nil {
		return 0
	}
	return len(keys)
}

//遍历缓存
func (this *redis) Each(call func(key, value interface{})) {
	keys, err := this.client.Keys(this.keyName + "*")
	if err != nil {
		return
	}
	for _, key := range keys {
		data, err := this.client.Get(key)
		if err != nil {
			continue
		}
		call(key, data)
	}
}
