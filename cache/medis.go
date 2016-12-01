package cache

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	SET   = "s"
	DEL   = "d"
	CLEAN = "c"
	DML   = []byte{19, 88, 03, 21}
)

//我的缓存
type Medis struct {
	Memory //继承
	file   string
}

func NewMedis(file string) *Medis {
	o := new(Medis)
	o.file = file
	o.readData()
	o.reSaveData()
	return o
}

//保存
func (this *Medis) Set(k string, v interface{}, exp ...int64) error {
	this.Memory.Set(k, v, exp...)
	return this.saveData(SET, k)
}

//删除
func (this *Medis) Del(k string) error {
	this.Memory.Del(k)
	return this.saveData(DEL, k)
}

//获取转换成v的格式，并将更新后的数据保存到内存里面
func (this *Medis) GetBy(key string, v interface{}) {
	this.Memory.GetBy(key, v)
	this.Memory.Set(key, v)
}

//清空
func (this *Medis) Clean() error {
	var err error
	oldData := this.Memory.data
	this.Memory.Clean()
	err = this.saveData(CLEAN, "")
	if err != nil {
		this.Memory.data = oldData
	}
	return err
}

//从文件里面读取数据
func (this *Medis) readData() {
	file, err := os.OpenFile(this.file, os.O_RDONLY, 0666)
	defer file.Close()
	if err != nil {
		return
	}
	buf := bufio.NewReader(file)
	for {
		data, e := buf.ReadBytes('\n')
		if e != nil {
			break
		}
		this.parseData(data[:len(data)-1])
	}
}

//解析数据
func (this *Medis) parseData(data []byte) error {
	if len(data) < 10 {
		return errors.New("to short.")
	}
	this.Lock()
	defer this.Unlock()
	for k, b := range data {
		data[k] = 255 - b
	}
	this.autoInit()
	btArr := bytes.Split(data, DML)
	btLen := len(btArr)
	tp := string(btArr[0])
	if tp == SET && btLen >= 3 {
		cdata := new(CacheData)
		e := json.Unmarshal(btArr[2], &cdata)
		if e != nil {
			return e
		}
		this.Memory.data[string(btArr[1])] = cdata
	} else if tp == DEL && btLen >= 2 {
		delete(this.Memory.data, string(btArr[1]))
	} else if tp == CLEAN {
		this.Memory.data = make(map[interface{}]*CacheData)
	}
	return nil
}

//更新保存(重新保存一次)
func (this *Medis) reSaveData() error {
	if this.Len() == 0 {
		return errors.New("it's null,can't write in file.")
	}
	err := ioutil.WriteFile(this.file, []byte{}, 0666)
	if err != nil {
		return err
	}
	this.Each(func(k, v interface{}) {
		if key, ok := k.(string); ok {
			this.Set(key, v)
		}
	})
	return nil
}

//保存数据到文件
//如果是del，则不需要保存data数据
func (this *Medis) saveData(cmd string, key string) (err error) {
	bt := []byte(cmd)
	bt = append(bt, DML...)
	bt = append(bt, []byte(key)...)
	bt = append(bt, DML...)
	if cmd == SET {
		d := this.Memory.get(key)
		data, err := json.Marshal(d)
		if err != nil {
			return err
		}
		bt = append(bt, data...)
	}
	for k, b := range bt {
		bt[k] = 255 - b
	}
	bt = append(bt, '\n')
	file, err := os.OpenFile(this.file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	defer file.Close()
	if err != nil {
		this.Memory.Del(key)
		return err
	}
	_, err = file.Write(bt)
	if err != nil {
		this.Memory.Del(key)
		return err
	}
	return nil
}
