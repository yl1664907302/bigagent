package utils

import (
	"bytes"
	"encoding/binary"
)

// 整体转换成字节
func InToBytes(n int) []byte {
	x := int64(n)
	bytesbuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesbuffer, binary.BigEndian, x)
	return bytesbuffer.Bytes()
}

// 字节转成整形
func BytesToInt(b []byte) int {
	bytesbuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesbuffer, binary.BigEndian, &x)
	return int(x)
}
