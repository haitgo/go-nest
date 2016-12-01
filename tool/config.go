package tool

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

//读取配置文件(自动替换配置文件里面的/**/注释)
//输入：file=配置文件名，obj=配置对象
func Config(file string, obj interface{}) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		log.Fatalln("配置文件读取失败：", err.Error())
	}
	compile, _ := regexp.Compile(`\/\*[^\*]+\*\/`)
	data = compile.ReplaceAll(data, []byte(""))
	err = json.Unmarshal(data, obj)
	if err != nil {
		log.Fatalln("配置文件解析错误：", err.Error())
	}
}
