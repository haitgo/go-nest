package cache

type Array struct {
	data   []interface{}
	length int
}

func NewArray(l ...int) *Array {
	o := new(Array)
	o.data = make([]interface{}, 0)
	if len(l) == 1 {
		o.length = l[0]
	}
	return o
}

//向数组的末尾添加一个或更多元素，并返回新的长度
func (this *Array) Push(d ...interface{}) int {
	if this.length > 0 && this.Length()+len(d) > this.length {
		this.Shift()
	}
	this.data = append(this.data, d...)
	return len(this.data)
}

//向数组的开头添加一个或更多元素，并返回新的长度。
func (this *Array) Unshift(d ...interface{}) int {
	if this.length > 0 && this.Length()+len(d) > this.length {
		this.Pop()
	}
	this.data = append(d, this.data...)
	return len(this.data)
}

//	删除并返回数组的第一个元素
func (this *Array) Shift() interface{} {
	d := this.data[0]
	this.data = this.data[1:]
	return d
}

//pop()	删除并返回数组的最后一个元素
func (this *Array) Pop() interface{} {
	l := len(this.data)
	if l == 0 {
		return nil
	}
	d := this.data[l-1]
	this.data = this.data[:l-1]
	return d
}

//从某个已有的数组返回选定的元素
func (this *Array) Slice(start, stop int) []interface{} {
	return this.data[start:stop]
}

//数组长度
func (this *Array) Length() int {
	return len(this.data)
}

//颠倒数组中元素的顺序。
func (this *Array) Reverse() {
	d := make([]interface{}, 0)
	l := len(this.data)
	for i := l; i >= 0; i-- {
		d = append(d, this.data[i])
	}
}
func (this *Array) Get(index int) interface{} {
	return this.data[index]
}

func (this *Array) Data() []interface{} {
	return this.data
}

/*
concat()	连接两个或更多的数组，并返回结果。
join()	把数组的所有元素放入一个字符串。元素通过指定的分隔符进行分隔。
pop()	删除并返回数组的最后一个元素
push()	向数组的末尾添加一个或更多元素，并返回新的长度。
reverse()	颠倒数组中元素的顺序。
shift()	删除并返回数组的第一个元素
slice()	从某个已有的数组返回选定的元素
sort()	对数组的元素进行排序
splice()	删除元素，并向数组添加新元素。
toSource()	返回该对象的源代码。
toString()	把数组转换为字符串，并返回结果。
toLocaleString()	把数组转换为本地数组，并返回结果。
unshift()	向数组的开头添加一个或更多元素，并返回新的长度。
valueOf()	返回数组对象的原始值
*/
