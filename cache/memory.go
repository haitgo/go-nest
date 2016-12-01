package cache

import (
	"encoding/json"
	"sync"
	"time"
)

//本缓存包，可以方便的使用缓存来保存热数据

//缓存数据结构
type CacheData struct {
	Value   interface{} //缓存数据的值
	Timeout int64       //到期时间(单位:秒）
}

//内存缓存
type Memory struct {
	data       map[interface{}]*CacheData //数据结构
	sync.Mutex                            //互斥锁
	Size       int                        //初始内存大小
	maked      bool                       //map是否已经初始化
	expire     int64                      //延迟时间
}

//创建一个内存缓存
//size=初始内存大小
func NewMemory(size int) *Memory {
	obj := new(Memory)
	obj.Size = size
	return obj
}

//导入数据
func (this *Memory) LoadData(data map[interface{}]interface{}) {
	this.Lock()
	defer this.Unlock()
	this.autoInit()
	for k, v := range data {
		vl := new(CacheData)
		vl.Value = v
		this.data[k] = vl
	}
}

//自动初始化map
func (this *Memory) autoInit() {
	if !this.maked {
		this.data = make(map[interface{}]*CacheData, this.Size)
		this.maked = true
	}
}

//检查是否到期，如果到期则自动除该数据
func (this *Memory) inspectExpire() {
	for key, data := range this.data {
		if this.isExpire(data) {
			delete(this.data, key)
		}
	}
}

//延时，（对所有的数据进行延时）
func (this *Memory) Expire(exp int64) {
	this.expire = exp
	this.Lock()
	defer this.Unlock()
	for _, data := range this.data {
		data.Timeout = time.Now().Unix() + exp
	}
}

//保存缓存
//key=键名，value=值，expire=生存时间，生存时间=0时永久存储
//如果当前的大小已经大于初始化的空间，则自动释放部分空间
func (this *Memory) Set(key, value interface{}, exp ...int64) error {
	this.autoInit()
	this.Lock()
	defer this.Unlock()
	if len(this.data) > this.Size {
		this.inspectExpire()
	}
	var cq *CacheData
	if d, ok := this.data[key]; ok {
		cq = d
	} else {
		cq = new(CacheData)
	}
	cq.Value = value
	expire := this.expire
	if len(exp) == 1 {
		expire = exp[0]
	}
	if expire != 0 {
		cq.Timeout = time.Now().Unix() + expire
	}
	this.data[key] = cq
	return nil
}

//读取缓存
//key=键名,返回值必须是未到期的，或者到期时间为0的
func (this *Memory) Get(key interface{}) interface{} {
	this.autoInit()
	this.Lock()
	defer this.Unlock()
	if data, ok := this.data[key]; ok && !this.isExpire(data) {
		return data.Value
	}
	delete(this.data, key)
	return nil
}
func (this *Memory) get(key interface{}) *CacheData {
	this.Lock()
	this.Unlock()
	return this.data[key]
}

//获取值为int类型
func (this *Memory) GetInt(key interface{}) int {
	if i, ok := this.Get(key).(int); ok {
		return i
	}
	return 0
}

//获取值为float类型
func (this *Memory) GetFloat(key interface{}) float64 {
	if i, ok := this.Get(key).(float64); ok {
		return i
	}
	return 0
}

//获取值为string类型
func (this *Memory) GetString(key interface{}) string {
	if i, ok := this.Get(key).(string); ok {
		return i
	}
	return ""
}

//获取值，并自动转换赋值给value
func (this *Memory) GetBy(key, value interface{}) {
	v := this.Get(key)
	if v == nil {
		return
	}
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	json.Unmarshal(data, value)
}

//删除缓存,如果key存在
func (this *Memory) Del(key interface{}) {
	this.autoInit()
	this.Lock()
	defer this.Unlock()
	delete(this.data, key)
}

//清空
func (this *Memory) Clean() {
	this.Lock()
	defer this.Unlock()
	this.data = make(map[interface{}]*CacheData, this.Size)
	this.maked = true
}

//判断是否到期，true到期，false未到期
func (this *Memory) isExpire(data *CacheData) bool {
	if data.Timeout == 0 || data.Timeout > time.Now().Unix() {
		return false
	}
	return true
}

//缓存数据数量
func (this *Memory) Len() int {
	this.autoInit()
	this.Lock()
	defer this.Unlock()
	this.inspectExpire()
	return len(this.data)
}

//迭代,只迭代未到期的缓存数据
func (this *Memory) Each(callback func(key, value interface{})) {
	this.autoInit()
	for key, data := range this.data { //
		if !this.isExpire(data) { //未到期
			callback(key, data.Value) //执行回调
		} else {
			delete(this.data, key)
		}
	}
}
