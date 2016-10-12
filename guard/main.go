package main

//守护进程工具，保证进程不会自动退出
//使用方法：guard  程序名 参数1 参数...
import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("故障退出", err)
			os.Exit(0)
		}
	}()
	if len(os.Args) < 2 {
		fmt.Println("没有输入要启动的程序")
		return
	}
	filePath, err := exec.LookPath(os.Args[1])
	if err != nil {
		fmt.Println("该命令不存在", err.Error())
		return
	}
	var stop int = 0
	for { //死循环，保证不会自己退出
		if stop > 10 {

		}
		arg := os.Args[2:]
		cmd := exec.Command(filePath, arg...)
		cmd.Stdin = os.Stdin   //将其他命令传入生成出的进程
		cmd.Stdout = os.Stdout //给新进程设置文件描述符，可以重定向到文件中
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println("-----", "子进程遇到错误退出", err.Error(), "-----")
		}
		stop++
	}
}
