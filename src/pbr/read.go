package smn_pbr

//product by tools, should not change this file.
//Author: SureMoon
//
import "github.com/golang/protobuf/proto"
import "pb/dict"
import "pb/base"
import "pb/rip_rpc_itf"
var funcList = []funcGetMsg {
    dict.EDict_base_Call:base_Call,
    dict.EDict_rip_rpc_itf_Login_DoLogin_Ret:rip_rpc_itf_Login_DoLogin_Ret,
    dict.EDict_rip_rpc_itf_Login_Test1_Ret:rip_rpc_itf_Login_Test1_Ret,
    dict.EDict_rip_rpc_itf_Login_DoLogin_Prm:rip_rpc_itf_Login_DoLogin_Prm,
    dict.EDict_rip_rpc_itf_Login_Test1_Prm:rip_rpc_itf_Login_Test1_Prm,
    dict.EDict_base_Ret:base_Ret,
}
func base_Call(bytes []byte) proto.Message {
    msg := &base.Call{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_DoLogin_Ret(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_DoLogin_Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_Test1_Ret(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test1_Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_DoLogin_Prm(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_DoLogin_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_Test1_Prm(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test1_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func base_Ret(bytes []byte) proto.Message {
    msg := &base.Ret{}
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


