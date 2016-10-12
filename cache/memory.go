package cache

import (
	"reflect"
	"sync"
	"time"
)

//本缓存包，可以方便的使用缓存来保存热数据

//缓存数据结构
type CacheData struct {
	value   interface{} //缓存数据的值
	timeout int64       //到期时间(单位:秒）
}

//内存缓存
type memory struct {
	data   map[interface{}]*CacheData //数据结构
	lock   *sync.Mutex                //互斥锁
	size   int64                      //初始内存大小
	expire int64                      //延迟时间
}

//创建一个内存缓存
//size=初始内存大小
func Memory(size int64) *memory {
	o := new(memory)
	o.data = make(map[interface{}]*CacheData, size)
	o.lock = new(sync.Mutex)
	o.size = size
	return o
}

//检查是否到期，如果到期则自动除该数据
func (this *memory) inspectExpire() {
	for key, data := range this.data {
		if this.isExpire(data) {
			delete(this.data, key)
		}
	}
}

//延时
func (this *memory) Expire(exp int64) {
	this.expire = exp
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, data := range this.data {
		data.timeout = time.Now().Unix() + exp
	}
}

//保存缓存
//key=键名，value=值，expire=生存时间，生存时间=0时永久存储
//如果当前的大小已经大于初始化的空间，则自动释放部分空间
func (this *memory) Set(key, value interface{}, exp ...int64) error {
	this.lock.Lock()
	defer this.lock.Unlock()
	length := len(this.data)
	if int64(length) > this.size {
		this.inspectExpire()
	}
	var cq *CacheData
	if d, ok := this.data[key]; ok {
		cq = d
	} else {
		cq = new(CacheData)
	}
	cq.value = value
	expire := this.expire
	if len(exp) == 1 {
		expire = exp[0]
	}
	if expire != 0 {
		cq.timeout = time.Now().Unix() + expire
	}
	this.data[key] = cq
	return nil
}

//读取缓存
//key=键名,返回值必须是未到期的，或者到期时间为0的
func (this *memory) Get(key interface{}) interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	if data, ok := this.data[key]; ok && !this.isExpire(data) {
		return data.value
	}
	delete(this.data, key)
	return nil
}

//获取值，并自动转换赋值给value
func (this *memory) GetBy(key, value interface{}) {
	v := this.Get(key)
	if v != nil {
		vf := reflect.ValueOf(value).Elem()
		vl := reflect.ValueOf(v)
		vf.Set(vl)
	}
}

//删除缓存,如果key存在
func (this *memory) Del(key interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.data, key)
}

//判断是否到期，true到期，false未到期
func (this *memory) isExpire(data *CacheData) bool {
	if data.timeout == 0 || data.timeout > time.Now().Unix() {
		return false
	}
	return true
}

//缓存数据数量
func (this *memory) Len() int {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.inspectExpire()
	return len(this.data)
}

//迭代,只迭代未到期的缓存数据
func (this *memory) Each(callback func(key, value interface{})) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for key, data := range this.data {
		if !this.isExpire(data) {
			callback(key, data.value)
		}
		delete(this.data, key)
	}
}
