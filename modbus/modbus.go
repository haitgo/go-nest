//modbus数据处理工具
package modbus

import (
	"go-nest/tool"
)

//编译采集命令
func Compile(addr, code int8, start, length int16) []byte {
	data := []byte{}
	data = append(data, tool.NumberToByte(addr)...)
	data = append(data, tool.NumberToByte(code)...)
	data = append(data, tool.NumberToByte(start)...)
	data = append(data, tool.NumberToByte(length)...)
	data = append(data, crc(data)...)
	return data
}

//解析数据
func ParseData(data []byte) *ModBus {
	//crc数据校验
	length := len(data) //数据长度
	if length < 5 {
		return nil
	}
	mb := new(ModBus)
	tool.ByteToNumber(data[:1], &mb.Addr)
	tool.ByteToNumber(data[1:2], &mb.Code)
	tool.ByteToNumber(data[2:3], &mb.Length)
	mb.Data = data[3 : length-2]
	mb.CRC = data[length-2:]
	if string(crc(data[:length-2])) != string(mb.CRC) {
		return nil
	}
	return mb
}

//modbus结构体
type ModBus struct {
	Addr   int8   //从机地址
	Code   int8   //功能码
	Length int8   //数据长度
	Data   []byte //数据
	CRC    []byte //crc校验
}

//crc校验码
func crc(data []byte) []byte {
	var crc16 uint16 = 0xffff
	l := len(data)
	for i := 0; i < l; i++ {
		crc16 ^= uint16(data[i])
		for j := 0; j < 8; j++ {
			if crc16&0x0001 > 0 {
				crc16 = (crc16 >> 1) ^ 0xA001
			} else {
				crc16 >>= 1
			}
		}
	}
	packet := make([]byte, 2)
	packet[0] = byte(crc16 & 0xff)
	packet[1] = byte(crc16 >> 8)
	return packet
}
