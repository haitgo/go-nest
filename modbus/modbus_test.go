package modbus

import (
	"testing"
)

func TestMobus(t *testing.T) {
	cmd := Compile(12, 3, 1, 3)
	t.Log("命令是：", cmd)
	d := []byte{12, 3, 5, 1, 2, 2, 1, 2}
	//d = append(d, crc(d)...)
	e := ParseData(d)
	t.Log("结果是", e)
}
