package proto_tool

import (
	"github.com/golang/protobuf/proto"
	"pb/base"
	"pb/dict"
)

var funcList = []funcGetMsg{
	dict.EDict_base_Call: base_Call,
}

func base_Call(bytes []byte) proto.Message {
	msg := &base.Call{}
	proto.Unmarshal(bytes, msg)
	return msg
}

type funcGetMsg func(bytes []byte) proto.Message

func GetMsgByDict(bytes []byte, dictId int) proto.Message {
	if dictId >= len(funcList) || dictId < 0 || funcList[dictId] == nil {
		return nil
	}
	return funcList[dictId](bytes)
}
