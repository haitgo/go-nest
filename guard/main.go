package main

//守护进程工具，保证进程不会自动退出
//使用方法：guard  程序名 参数1 参数...
//如果在10秒内启动次数超过10次表示守护失败
import (
	"flag"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	logs       *log.Logger     //日志
	reloadNums int64       = 0 //重启次数
)

func main() {
	t := flag.Int64("t", 10, "守护检测时间,单位:秒")
	n := flag.Int64("n", 10, "检测时间内允许重启的次数")
	cmd := flag.String("c", "", "守护执行的命令,例如\"abc -a=12 -c\"")
	flag.Parse()
	file, _ := os.OpenFile("guard.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 666)
	defer func() {
		file.Close()
		if err := recover(); err != nil {
			log.Fatal("故障退出", err)
		}
	}()
	if *cmd == "" {
		log.Fatal("请输入要启动的程序")
	}
	arg := strings.Split(*cmd, " ")
	filePath, err := exec.LookPath(arg[0])
	if err != nil {
		log.Fatal("该命令不存在", err.Error())
	}
	logs = log.New(file, filePath+"\t", 1|2|3) //创建日志
	var lastTime int64 = 0
	for { //死循环，保证不会自己退出
		nowTime := time.Now().Unix()
		if lastTime == 0 {
			lastTime = nowTime
		}
		if nowTime-lastTime <= *t && reloadNums >= *n {
			log.Println(filePath, "守护失败")
			logs.Fatal("守护失败")
		}
		reloadNums++
		arga := arg[1:]
		cmd := exec.Command(filePath, arga...)
		cmd.Stdin = os.Stdin   //将其他命令传入生成出的进程
		cmd.Stdout = os.Stdout //给新进程设置文件描述符，可以重定向到文件中
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil { //
			logs.Println("遇到错误退出", err.Error())
		}
		go func() {
			time.Sleep(time.Second * 12)
			reloadNums = 0
			lastTime = 0
		}()
	}
}
