package main

import (
	"bytes"
	"com.suremoon.net/basis/smn_net"
	"fmt"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	buf := bytes.NewBuffer(make([]byte, 0))
	str := "hello world"
	err := smn_net.WriteString(str, buf)
	checkerr(err)
	fmt.Println(buf.Bytes())
	r_str, err := smn_net.ReadString(buf)
	fmt.Println(r_str)
	buf.Reset()
	err = smn_net.WriteInt(132, buf)
	checkerr(err)
	fmt.Println(buf.Bytes())
	i, err := smn_net.ReadInt(buf)
	fmt.Println(i)

}
