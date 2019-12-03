package smn_rpc

import (
	"net"
	"time"
	"com.suremoon.net/basis/smn_stream"
	"com.suremoon.net/smn/net_libs/muti_service"
	"pb/smn_base"
	"com.suremoon.net/basis/smn_net"
	"github.com/golang/protobuf/proto"
)

func NewRPCFConn(no int64, mgr *muti_service.ServiceManager, desc string, localAddr, remoteAddr net.Addr) muti_service.ForwardConnItf {
	return  &rpcFConn{no: no, mgr: mgr, desc: desc, localAddr:localAddr, remoteAddr:remoteAddr, cache:smn_stream.NewByteCache(1000, 1 * time.Second), status:STATUS_READY}
}

func NewRPCServiceManager(conn net.Conn) *muti_service.ServiceManager {
	sm := muti_service.NewServiceManager(conn)
	sm.FConnFactory = NewRPCFConn
	return sm
}

func ServiceManagerRegister(mgr *muti_service.ServiceManager, no int64, desc string, rpcSvr RpcSvrItf) (conn muti_service.ForwardConnItf, isExist bool) {
	conn, isExist = mgr.Regitster(no, desc)
	if isExist{
		return
	}
	rpcFc := conn.(*rpcFConn)
	rpcFc.rpcSvr = rpcSvr
	return rpcFc, isExist
}

type RPCStatus int

const (
	STATUS_READY RPCStatus = iota
	STATUS_READLEN
)

type rpcFConn struct {
	localAddr  net.Addr
	remoteAddr net.Addr
	no         int64
	mgr        *muti_service.ServiceManager
	desc       string
	cache      *smn_stream.ByteCache
	status     RPCStatus
	callLen    int
	rpcSvr     RpcSvrItf
}

func (this *rpcFConn) Read(b []byte) (n int, err error) {
	return this.cache.Read(b)
}

func (this *rpcFConn) Write(b []byte) (n int, err error) {
	this.SendToSM(b)
	return 0, nil
}

func (this *rpcFConn) Close() error {
	this.mgr.Drop(this.no)
	return nil
}

func (this *rpcFConn) LocalAddr() net.Addr {
	return this.localAddr
}

func (this *rpcFConn) RemoteAddr() net.Addr {
	return this.remoteAddr
}

func (this *rpcFConn) SetDeadline(t time.Time) error {
	return nil
}

func (this *rpcFConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (this *rpcFConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (this *rpcFConn) SendToSM(msg []byte) {
	this.mgr.Send(&smn_base.FPkg{NO: this.no, Msg: msg})
}

func (this *rpcFConn) RecvFromSM(msg []byte) {
	this.cache.Write(msg)
	switch this.status {
	case STATUS_READY:{
		if this.cache.Len() >= 8{
			this.callLen, _ = smn_net.ReadInt(this)
			this.status = STATUS_READLEN
		}
	}
	case STATUS_READLEN:
		if this.cache.Len() >= this.callLen{
			bts := make([]byte, this.callLen)
			this.Read(bts)
			callMsg := &smn_base.Call{}
			proto.Unmarshal(bts, callMsg)
			dict, res, err := this.rpcSvr.OnMessage(callMsg, this)
			WriteRet(this, dict, res, err)
			this.status = STATUS_READY
		}

	}
}

func (this *rpcFConn) Desc() string {
	return this.desc
}

func (this *rpcFConn) SetTimeOut(t time.Duration){
	this.cache.TimeOut = t
}
