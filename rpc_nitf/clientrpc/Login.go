package clt_rpc_rpc_itf

//Product by SureMoonNet
//Author: ProtossGenius
//Auto-code should not change.
import "github.com/ProtossGenius/SureMoonNet/test/rpc_itf"
import "github.com/golang/protobuf/proto"
import "github.com/ProtossGenius/SureMoonNet/pb/rip_rpc_itf"
import "github.com/ProtossGenius/SureMoonNet/pb/smn_dict"
import "github.com/ProtossGenius/SureMoonNet/smn/net_libs/smn_rpc"
import "sync"
type CltRpcLogin struct {
    rpc_itf.Login
    conn smn_rpc.MessageAdapterItf
    lock sync.Mutex
}
func NewCltRpcLogin(conn smn_rpc.MessageAdapterItf) *CltRpcLogin {
    return &CltRpcLogin{conn:conn}
}
func (this *CltRpcLogin)DoLogin(user string, pswd string, code int) (bool, int) {
    this.lock.Lock()
    defer this.lock.Unlock()
    msg := &rip_rpc_itf.Login_DoLogin_Prm{User:user, Pswd:pswd, Code:int64(code)}
    this.conn.WriteCall(smn_dict.EDict_rip_rpc_itf_Login_DoLogin_Prm, msg)
    rm, err := this.conn.ReadRet()
    if err != nil{
    	panic(err)
    }
    if rm.Err{
    	panic(string(rm.Msg))
    }
    res := &rip_rpc_itf.Login_DoLogin_Ret{}
    err = proto.Unmarshal(rm.Msg, res)
    if err != nil{
    	panic(err)
    }
    return res.P0, int(res.P1)
}
func (this *CltRpcLogin)Test1(a []string, b []int, c []uint, d []uint64, e []int32) ([]int) {
    this.lock.Lock()
    defer this.lock.Unlock()
    msg := &rip_rpc_itf.Login_Test1_Prm{A:a, B:smn_rpc.IntArrToInt64Arr(b), C:smn_rpc.UIntArrToUInt64Arr(c), D:d, E:e}
    this.conn.WriteCall(smn_dict.EDict_rip_rpc_itf_Login_Test1_Prm, msg)
    rm, err := this.conn.ReadRet()
    if err != nil{
    	panic(err)
    }
    if rm.Err{
    	panic(string(rm.Msg))
    }
    res := &rip_rpc_itf.Login_Test1_Ret{}
    err = proto.Unmarshal(rm.Msg, res)
    if err != nil{
    	panic(err)
    }
    return smn_rpc.Int64ArrToIntArr(res.P0)
}
func (this *CltRpcLogin)Test2(key string, c smn_rpc.ConnFunc) (bool) {
    this.lock.Lock()
    defer this.lock.Unlock()
    msg := &rip_rpc_itf.Login_Test2_Prm{Key:key}
    this.conn.WriteCall(smn_dict.EDict_rip_rpc_itf_Login_Test2_Prm, msg)
    c(this.conn.GetConn())
    rm, err := this.conn.ReadRet()
    if err != nil{
    	panic(err)
    }
    if rm.Err{
    	panic(string(rm.Msg))
    }
    res := &rip_rpc_itf.Login_Test2_Ret{}
    err = proto.Unmarshal(rm.Msg, res)
    if err != nil{
    	panic(err)
    }
    return res.P0
}
