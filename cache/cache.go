package cache

//缓存接口，所有的缓存必须满足该接口
type Cacher interface {
	Expire(exp int64)
	Set(key, value interface{}, expire ...int64) error
	Get(key interface{}) interface{}
	GetBy(key, value interface{})
	Del(key interface{})
	Len() int
	Each(func(k, v interface{}))
}

//缓存类
type Cache struct {
	cacher Cacher
}

func New(cacher Cacher) *Cache {
	o := new(Cache)
	o.cacher = cacher
	return o
}

//延时
func (this *Cache) Expire(exp int64) *Cache {
	this.cacher.Expire(exp)
	return this
}

//保存值
//expire=生存时间，超时将自动删除，单位：秒,可以为空，如果为空，则使用Expire()设置的时间
func (this *Cache) Set(key, value interface{}, expire ...int64) error {
	return this.cacher.Set(key, value, expire...)
}

//默认获取值的方法
func (this *Cache) Get(key interface{}) interface{} {
	return this.cacher.Get(key)
}

//获取值为int类型
func (this *Cache) GetInt(key interface{}) int64 {
	if i, ok := this.cacher.Get(key).(int64); ok {
		return i
	}
	return 0
}

//获取值为float类型
func (this *Cache) GetFloat(key interface{}) float64 {
	if i, ok := this.cacher.Get(key).(float64); ok {
		return i
	}
	return 0
}

//获取值为string类型
func (this *Cache) GetString(key interface{}) string {
	if i, ok := this.cacher.Get(key).(string); ok {
		return i
	}
	return ""
}

//获取值，并自动转换赋值给value
func (this *Cache) GetBy(key interface{}, value interface{}) {
	this.cacher.GetBy(key, value)
}

//获取缓存的数据的数量
func (this *Cache) Len() int {
	return this.cacher.Len()
}

//删除这个缓存
func (this *Cache) Del(key interface{}) {
	this.cacher.Del(key)
}

//遍历所有的缓存数据，并回调
func (this *Cache) Each(call func(key, value interface{})) {
	this.cacher.Each(call)
}
