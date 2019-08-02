package smn_pbr

//product by tools, should not change this file.
//Author: SureMoon
//
import "github.com/golang/protobuf/proto"
import "pb/dict"
import "pb/base"
import "pb/rpc_itf"
var funcList = []funcGetMsg {
    dict.EDict_base_Call:base_Call,
    dict.EDict_base_Qnm:base_Qnm,
    dict.EDict_rpc_itf_Rpc_Itf_DoLogin_Prm:rpc_itf_Rpc_Itf_DoLogin_Prm,
    dict.EDict_rpc_itf_Rpc_Itf_DoLogin_Ret:rpc_itf_Rpc_Itf_DoLogin_Ret,
}
func base_Call(bytes []byte) proto.Message {
    msg := &base.Call{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func base_Qnm(bytes []byte) proto.Message {
    msg := &base.Qnm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rpc_itf_Rpc_Itf_DoLogin_Prm(bytes []byte) proto.Message {
    msg := &rpc_itf.Rpc_Itf_DoLogin_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rpc_itf_Rpc_Itf_DoLogin_Ret(bytes []byte) proto.Message {
    msg := &rpc_itf.Rpc_Itf_DoLogin_Ret{}
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


