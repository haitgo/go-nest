//bind仅适用于web开发使用，可以绑定post和get的参数到结构体
package bind

import (
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//bind 结构体
type Bind struct {
	Data    map[string]string //绑定的数据源
	needArr map[string]string //必须要的
	omitArr map[string]string //排除的
}

//绑定map
func Map(data map[string]string) *Bind {
	o := new(Bind)
	o.Data = data
	o.needArr = make(map[string]string)
	o.omitArr = make(map[string]string)
	return o
}

//绑定表单
func Form(req *http.Request) *Bind {
	o := new(Bind)
	o.needArr = make(map[string]string)
	o.omitArr = make(map[string]string)
	req.ParseMultipartForm(32 << 20) // 32 MB
	data := map[string]string{}
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		for k, v := range req.MultipartForm.Value {
			data[k] = strings.TrimSpace(v[0])
		}
	} else {
		for k, v := range req.PostForm {
			data[k] = strings.TrimSpace(v[0])
		}
	}
	o.Data = data
	return o
}

//绑定url查询属性
func Query(req *http.Request) *Bind {
	o := new(Bind)
	o.needArr = make(map[string]string)
	o.omitArr = make(map[string]string)
	data := map[string]string{}
	for k, v := range req.URL.Query() {
		data[k] = strings.TrimSpace(v[0])
	}
	o.Data = data
	return o
}

//map必须包含arr的每个元素
func (this *Bind) Need(arr ...string) *Bind {
	for _, v := range arr {
		this.needArr[v] = v
	}
	return this
}

//map省略arr的每个元素，当出现该参数的时候自动忽略
func (this *Bind) Omit(arr ...string) *Bind {
	for _, v := range arr {
		this.omitArr[v] = v
	}
	return this
}

//赋值
func (this *Bind) Set(obj interface{}) (err error, skey string) {
	ref := reflect.ValueOf(obj).Elem()
	typ := ref.Type()
	return this.refSet(ref, typ)
}
func (this *Bind) refSet(ref reflect.Value, typ reflect.Type) (err error, skey string) {
	for i := 0; i < ref.NumField(); i++ {
		if ref.Field(i).Kind() == reflect.Struct { //如果为结构体，则循环从结构体里面去绑定参数
			er, ek := this.refSet(ref.Field(i), typ.Field(i).Type)
			if er != nil {
				return er, ek
			}
			continue
		}
		field := typ.Field(i)
		valid := field.Tag.Get("valid") //获取到的验证参数
		if valid == "" {                //如果不开启验证的
			continue
		}
		def := field.Tag.Get("default")        //默认值
		skey = field.Name                      //结构体成员属性名称
		if _, has := this.omitArr[skey]; has { //如果此key的参数是需要省略的则不过滤
			continue
		}

		value := this.Data[skey] //根据结构体名获取到map对应的值
		mustFilter := true
		if _, has := this.needArr[skey]; has { //当必须包含的时候
			delete(this.needArr, skey)
		} else {
			mustFilter = false
		}
		if value == "" { //当value为空的时候就有3种情况
			if def != "" { //如果有默认值的时候则赋默认值
				value = def
			} else if mustFilter { //如果是必须过滤的则不能为空
				err = errors.New("不能为空")
				return
			} else { //否则就不过滤
				continue
			}

		}
		if !ref.Field(i).CanSet() { //如果该字段不允许赋值
			err = errors.New("拒绝写入")
			return
		}
		if !this.regexp(value, valid) { //如果正则判断不正确
			err = errors.New("数据异常")
			return
		}
		//类型断言转换
		rvalue, errb := this.typeConversion(value, field.Type.Name())
		if errb != nil { //如果断言转换类型出现错误
			err = errb
			return
		}
		ref.Field(i).Set(rvalue) //给结构体赋值
	}
	if len(this.needArr) > 0 { //当必须包含的时候但是缺乏元素
		err = errors.New("不能为空")
		for k, _ := range this.needArr {
			skey = k //返回第一个未包含的skey
			continue
		}
	}
	return
}

//类型断言
func (this *Bind) typeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "uint" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(uint(i)), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int16" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int16(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int32(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}
	//else if .......增加其他一些类型的转换
	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}

//
//格式校验 str 输入的字符串，partten正则式或快捷正则
func (this *Bind) regexp(str, parttens string) bool {
	parttenArr := strings.Split(parttens, ",")
	for _, partten := range parttenArr {
		partten = strings.TrimSpace(partten)
		if strings.Index(partten, "length") >= 0 { //如果有查找长度
			partten = strings.TrimLeft(partten, "length(")
			partten = strings.Replace(partten, "|", ",", -1)
			partten = strings.TrimRight(partten, ")")
			partten = `^[\S]{` + partten + `}$`
		} else if strings.Index(partten, "match") >= 0 { //如果为正则匹配
			partten = strings.TrimLeft(partten, "match(")
			partten = strings.TrimRight(partten, ")")
		}
		if !Regexp(str, partten) {
			return false
		}
	}
	return true
}
