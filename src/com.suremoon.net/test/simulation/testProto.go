package main

import (
	"com.suremoon.net/basis/smn_str"
	"fmt"
	"github.com/golang/protobuf/proto"
	"pb/base"
	"pb/dict"
	"pbr"
)

type itf interface {
	f(a int,
		b int) (int,
		int)
}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	fmt.Println(smn_str.InitialsUpper(""))
	c := &base.Call{Dict: 123}
	bytes, err := proto.Marshal(c)
	checkerr(err)
	msg := smn_pbr.GetMsgByDict(bytes, dict.EDict_base_Call)
	if msg != nil {
		fmt.Println(msg.(*base.Call).Dict)
	}
}
