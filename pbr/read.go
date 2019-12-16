package smn_pbr

//product by tools, should not change this file.
//Author: SureMoon
//
import "github.com/golang/protobuf/proto"
import "github.com/ProtossGenius/SureMoonNet/pb/smn_dict"
import "github.com/ProtossGenius/SureMoonNet/pb/rip_rpc_itf"
import "github.com/ProtossGenius/SureMoonNet/pb/smn_base"
var funcList = []funcGetMsg {
    smn_dict.EDict_rip_rpc_itf_Login_Test1_Prm:rip_rpc_itf_Login_Test1_Prm,
    smn_dict.EDict_rip_rpc_itf_Login_Test1_Ret:rip_rpc_itf_Login_Test1_Ret,
    smn_dict.EDict_rip_rpc_itf_Login_Test2_Prm:rip_rpc_itf_Login_Test2_Prm,
    smn_dict.EDict_rip_rpc_itf_Login_Test2_Ret:rip_rpc_itf_Login_Test2_Ret,
    smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Prm:rip_rpc_itf_Login_DoLogin_Prm,
    smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Ret:rip_rpc_itf_Login_DoLogin_Ret,
    smn_dict.EDict_smn_base_Call:smn_base_Call,
    smn_dict.EDict_smn_base_Ret:smn_base_Ret,
    smn_dict.EDict_smn_base_FPkg:smn_base_FPkg,
}
func rip_rpc_itf_Login_Test1_Prm(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test1_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_Test1_Ret(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test1_Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_Test2_Prm(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test2_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_Test2_Ret(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_Test2_Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_DoLogin_Prm(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_DoLogin_Prm{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func rip_rpc_itf_Login_DoLogin_Ret(bytes []byte) proto.Message {
    msg := &rip_rpc_itf.Login_DoLogin_Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func smn_base_Call(bytes []byte) proto.Message {
    msg := &smn_base.Call{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func smn_base_Ret(bytes []byte) proto.Message {
    msg := &smn_base.Ret{}
    proto.Unmarshal(bytes, msg)
    return msg
}
func smn_base_FPkg(bytes []byte) proto.Message {
    msg := &smn_base.FPkg{}
    proto.Unmarshal(bytes, msg)
    return msg
}
type funcGetMsg func(bytes []byte) proto.Message
func GetMsgByDict(bytes []byte, dict smn_dict.EDict) proto.Message {
	dictId := int(dict)
	if dictId >= len(funcList) || dictId < 0 || funcList[dictId] == nil {
		return nil
	}
	return funcList[dictId](bytes)
}


