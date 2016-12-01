package cache

//缓存接口，所有的缓存必须满足该接口
type Cacher interface {
	Expire(exp int64)                                  //延时
	Set(key, value interface{}, expire ...int64) error //添加数据
	Get(key interface{}) interface{}                   //获取数据
	GetInt(key interface{}) int                        //获取值，并转换成int类型
	GetString(key interface{}) string                  //获取值，并转换为string
	GetFloat(key interface{}) float64                  //获取值，并转换成float
	GetBy(key, value interface{})                      //获取数据，并自动转换类型
	Del(key interface{})                               //删除数据
	Len() int                                          //数据长度
	Each(func(k, v interface{}))                       //循环处理每条数据
	Clean()                                            //清空缓存
}
