package smn_pbr

//product by tools, should not change this file.
//Author: SureMoon
//
import "github.com/golang/protobuf/proto"
import "pb/dict"
import "pb/base"
var funcList = []funcGetMsg {
    dict.EDict_base_Call:base_Call,
}
func base_Call(bytes []byte) proto.Message {
    msg := &base.Call{}
    proto.Unmarshal(bytes, msg)
    return msg
}
type funcGetMsg func(bytes []byte) proto.Message
func GetMsgByDict(bytes []byte, dict dict.EDict) proto.Message {
	dictId := int(dict)
	if dictId >= len(funcList) || dictId < 0 || funcList[dictId] == nil {
		return nil
	}
	return funcList[dictId](bytes)
}


