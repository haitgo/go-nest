package tool

import (
	"testing"
)

func TestConfig(t *testing.T) {
	m := make(map[string]string)
	Config("config.json", &m)
	t.Log("读取的配置文件是：", m)
}
