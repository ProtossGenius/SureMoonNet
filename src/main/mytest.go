package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var i1 int64 = 511 // [00000000 00000000 ... 00000001 11111111] = [0 0 0 0 0 0 1 255]

	s1 := make([]byte, 0)
	buf := bytes.NewBuffer(s1)

	// 数字转 []byte, 网络字节序为大端字节序
	binary.Write(buf, binary.BigEndian, &i1)
	s1 = buf.Bytes()
	fmt.Println(s1)
}
