package svr_rpc_rpc_itf

//Product by SureMoonNet
//Author: ProtossGenius
//Auto-code should not change.
import "github.com/ProtossGenius/SureMoonNet/test/rpc_itf"
import "github.com/golang/protobuf/proto"
import "github.com/ProtossGenius/SureMoonNet/pb/smn_dict"
import "github.com/ProtossGenius/SureMoonNet/pb/smn_base"
import "github.com/ProtossGenius/SureMoonNet/pbr"
import "net"
import "fmt"
import "github.com/ProtossGenius/SureMoonNet/pb/rip_rpc_itf"
import "github.com/ProtossGenius/SureMoonNet/smn/net_libs/smn_rpc"
type SvrRpcLogin struct {
    itf rpc_itf.Login
    dicts []smn_dict.EDict
}
func NewSvrRpcLogin(itf rpc_itf.Login) *SvrRpcLogin {
    list := make([]smn_dict.EDict, 0)
    list = append(list, smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Prm)
    list = append(list, smn_dict.EDict_rip_rpc_itf_Login_Test1_Prm)
    list = append(list, smn_dict.EDict_rip_rpc_itf_Login_Test2_Prm)
    return &SvrRpcLogin{itf:itf, dicts:list}
}
func (this *SvrRpcLogin)getEDictList() []smn_dict.EDict {
    return this.dicts
}
func (this *SvrRpcLogin)OnMessage(c *smn_base.Call, conn net.Conn) (_d smn_dict.EDict, _p proto.Message, _e error) {
    defer func() {
    if err := recover(); err != nil {
            _p = nil
            _e = fmt.Errorf("%v", err)
    }
    }()
    m := smn_pbr.GetMsgByDict(c.Msg, c.Dict)
    switch c.Dict {
        case smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Prm: {
            _d = smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Ret
            msg := m.(*rip_rpc_itf.Login_DoLogin_Prm)
            p0, p1 := this.itf.DoLogin(msg.User, msg.Pswd, int(msg.Code))
            return _d, &rip_rpc_itf.Login_DoLogin_Ret{P0:p0, P1:int64(p1)            }, nil
        }
        case smn_dict.EDict_rip_rpc_itf_Login_Test1_Prm: {
            _d = smn_dict.EDict_rip_rpc_itf_Login_Test1_Ret
            msg := m.(*rip_rpc_itf.Login_Test1_Prm)
            p0 := this.itf.Test1(msg.A, smn_rpc.Int64ArrToIntArr(msg.B), smn_rpc.UInt64ArrToUIntArr(msg.C), msg.D, msg.E)
            return _d, &rip_rpc_itf.Login_Test1_Ret{P0:smn_rpc.IntArrToInt64Arr(p0)            }, nil
        }
        case smn_dict.EDict_rip_rpc_itf_Login_Test2_Prm: {
            _d = smn_dict.EDict_rip_rpc_itf_Login_Test2_Ret
            msg := m.(*rip_rpc_itf.Login_Test2_Prm)
            p0 := this.itf.Test2(msg.Key, conn)
            return _d, &rip_rpc_itf.Login_Test2_Ret{P0:p0            }, nil
        }
    }
    return -1, nil, nil
}
