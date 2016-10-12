package xlog

import (
	"fmt"
	"io"
	"log"
	"time"
)

func New(name string, w io.Writer) *Xlog {
	o := new(Xlog)
	o.logger = log.New(w, "["+name+"]\t", log.Ldate|log.Ltime)
	return o
}

//实例化一个调试状态
func Debug(name string) *Xlog {
	o := new(Xlog)
	o.debug = true
	o.name = name
	return o
}

type Xlog struct {
	logger *log.Logger
	name   string //日志的文件名
	debug  bool   //是否为调试
}

func (this *Xlog) Print(err ...interface{}) {
	if this.debug == true {
		tm := time.Now().Format("2006-01-02 15:04:05")
		fmt.Println("["+this.name+"]", tm, err)
		return
	}
	this.logger.Println(err...)
}
