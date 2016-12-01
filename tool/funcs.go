package tool

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

//md5加密
func Md5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// 产生随机字符串
//输入：size=字符串长度，kind=字符串类型 0数字，1小写字母，2大写字母，3数字大小写字母
func Rand(size int, kind int) string {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

//求数组最大值
func Max(args ...float64) float64 {
	max := args[0]
	for _, v := range args {
		if v > max {
			max = v
		}
	}
	return max
}

//求数组最小值
func Min(args ...float64) float64 {
	min := args[0]
	for _, v := range args {
		if v < min {
			min = v
		}
	}
	return min
}

//求数组和
func Sum(args ...float64) float64 {
	sum := args[0]
	for _, v := range args {
		sum = sum + v
	}
	return sum
}

//将byte类型转换成字符串，例如[]byte{1,2,3,4,5}转换成"1,2,3,4,5"
//输入：bt=字节码,dilimiter=连接字符串
func ByteToStr(bt []byte, dilimiter string) string {
	str := ""
	for _, v := range bt {
		str = str + strconv.Itoa(int(v)) + dilimiter
	}
	return str
}

//将byte转换成数字类型
//输入：b=byte,data=要转换输出的变量'请输入指针
func ByteToNumber(b []byte, data interface{}) {
	b_buf := bytes.NewBuffer(b)
	binary.Read(b_buf, binary.BigEndian, data)
	return
}

//数字转换成byte类型
func NumberToByte(data interface{}) []byte {
	b_buf := bytes.NewBuffer([]byte{})
	binary.Write(b_buf, binary.BigEndian, data)
	return b_buf.Bytes()
}

//打开一个文件，如果不存在则创建
func OpenFile(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
}

//当前时间
func Now(str ...string) string {
	format := "2006-01-02 15:04:05"
	if len(str) == 1 {
		format = str[0]
	}
	return time.Now().Format(format)
}

//将map快速转换为结构体
func MapToStruct(m interface{}, s interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, s)
}

//输出到日志
var logLock = new(sync.Mutex)

func SaveLog(logfile string, arg ...interface{}) error {
	logLock.Lock()
	f, e := os.OpenFile(logfile+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if e != nil {
		return e
	}
	defer func() {
		logLock.Unlock()
		f.Close()
	}()
	l := log.New(f, "", 0)
	a := []interface{}{}
	a = append(a, time.Now().Format("[2006-01-02 15:04:05]:"))
	a = append(a, arg...)
	l.Println(a...)
	return nil
}
