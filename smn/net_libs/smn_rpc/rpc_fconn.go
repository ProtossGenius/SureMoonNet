package smn_rpc

import (
	"net"
	"github.com/ProtossGenius/SureMoonNet/smn/net_libs/muti_service"
)


func NewRPCServiceManager(conn net.Conn) *muti_service.ServiceManager {
	sm := muti_service.NewServiceManager(conn)
	return sm
}

func AccpterRun(adapter MessageAdapterItf, rpcSvr RpcSvrItf) {
	for {
		msg, err := adapter.ReadCall()
		dict, res, err := rpcSvr.OnMessage(msg, adapter.GetConn())
		adapter.WriteRet(dict, res, err)
	}
}

func ServiceManagerRegister(mgr *muti_service.ServiceManager, no int64, desc string, rpcSvr RpcSvrItf) (conn muti_service.ForwardConnItf, isExist bool) {
	conn, isExist = mgr.Regitster(no, desc)
	conn.SetTimeOut(0)
	if isExist {
		return
	}
	go AccpterRun(NewMessageAdapter(conn), rpcSvr)
	return conn, isExist
}