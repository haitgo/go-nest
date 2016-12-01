package tool

import (
	"go-nest/bind"

	"github.com/gin-gonic/gin"
)

const (
	ECHO_JSON int8 = 1
	ECHO_HTML      = 2
)

//对gin开发json接口快速处理
type FastBind struct {
	Content     *gin.Context
	binding     []*bind.Bind  //绑定方法
	obj         []interface{} //要绑定的对象
	send        bool          //是否已发送过json
	echoType    int8          //输出类型
	errHTMLfile string        //错误的html模板路径
}

//实例化一个对象
func NewFastBind(c *gin.Context) *FastBind {
	obj := new(FastBind)
	obj.Content = c
	obj.echoType = ECHO_JSON
	return obj
}

//设置输出为html
func (this *FastBind) ErrHTML(file string) {
	this.echoType = ECHO_HTML
	this.errHTMLfile = file
}

//从url绑定参数,可以执行多次
func (this *FastBind) QueryTo(obj interface{}) *bind.Bind {
	binding := bind.Query(this.Content.Request)
	this.obj = append(this.obj, obj)
	this.binding = append(this.binding, binding)
	return binding
}

//从from绑定参数，可以执行多次
func (this *FastBind) FormTo(obj interface{}) *bind.Bind {
	binding := bind.Form(this.Content.Request)
	this.obj = append(this.obj, obj)
	this.binding = append(this.binding, binding)
	return binding
}

//当绑定完成后执行的回调函数
func (this *FastBind) Success(call func()) {
	for k, binding := range this.binding {
		err, name := binding.Set(this.obj[k])
		if err != nil {
			this.Ret(2, gin.H{"name": name, "msg": err.Error()})
			return
		}
	}
	call()
	if !this.send { //执行错误
		this.Ret(0, gin.H{"msg": "Server error"})
	}
}

//返回json数据
//输入：第一个参数=code,第二个参数=gin.H
func (this *FastBind) Ret(args ...interface{}) {
	if this.send {
		return
	}
	code := args[0].(int)
	data := gin.H{}
	if len(args) == 2 {
		data = args[1].(gin.H)
	}
	data["code"] = code
	//针对错误的处理
	if this.echoType == ECHO_HTML && this.errHTMLfile != "" {
		this.Content.HTML(200, this.errHTMLfile, data)
	} else {
		this.Content.JSON(200, data)
	}
	this.send = true
}
