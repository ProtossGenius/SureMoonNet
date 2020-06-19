package muti_service

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_net"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_stream"
	"github.com/ProtossGenius/SureMoonNet/pb/smn_base"
	"google.golang.org/protobuf/proto"
	"github.com/pkg/errors"
)

const (
	//ErrServerNumNotExist don's exist that server num.
	ErrServerNumNotExist = "ErrServerNumNotExist :[%d]"
)

//FConnFactory .
type FConnFactory func(no int64, mgr *ServiceManager, desc string, localAddr, remoteAddr net.Addr) ForwardConnItf

func newDftFConn(no int64, mgr *ServiceManager, desc string, localAddr, remoteAddr net.Addr) ForwardConnItf {
	return &FConn{no: no, mgr: mgr, desc: desc, localAddr: localAddr, remoteAddr: remoteAddr,
		cache: smn_stream.NewByteCache(1000, 1*time.Second)}
}

//ServiceManager .
type ServiceManager struct {
	OnErr        smn_err.OnErr
	TimeOut      time.Duration
	conn         net.Conn
	sendChan     chan *smn_base.FPkg
	regMap       map[int64]ForwardConnItf //key always > 0
	close        chan int
	FConnFactory FConnFactory
	mapLock      sync.Mutex //因为map操作并不频繁，所以不用计较锁的开销
}

func NewServiceManager(conn net.Conn) *ServiceManager {
	return &ServiceManager{conn: conn, close: make(chan int, 1), regMap: make(map[int64]ForwardConnItf),
		sendChan: make(chan *smn_base.FPkg, 1024), FConnFactory: newDftFConn, OnErr: smn_err.DftOnErr,
		TimeOut: 1 * time.Second}
}

func (this *ServiceManager) Send(p *smn_base.FPkg) {
	this.sendChan <- p
}

func (this *ServiceManager) recv(p *smn_base.FPkg) {
	conn, ok := this.GetFConn(p.NO)
	if ok {
		if p.Err {
			err := errors.New(string(p.Msg))
			conn.ErrClose(err)
			this.OnErr(err)
			return
		}
		conn.RecvFromSM(p.Msg)
	} else {
		err := fmt.Errorf(ErrServerNumNotExist, p.NO)
		this.OnErr(err)
		this.Send(&smn_base.FPkg{NO: p.NO, Err: true, Msg: []byte(err.Error())})
	}
}

func (this *ServiceManager) Regitster(no int64, desc string) (conn ForwardConnItf, isExist bool) {
	this.mapLock.Lock()
	defer this.mapLock.Unlock()
	if c, ok := this.regMap[no]; ok {
		return c, true
	}
	conn = this.FConnFactory(no, this, desc, this.conn.LocalAddr(), this.conn.RemoteAddr())
	this.regMap[no] = conn
	return conn, false
}

func (this *ServiceManager) GetFConn(no int64) (ForwardConnItf, bool) {
	this.mapLock.Lock()
	defer this.mapLock.Unlock()
	c, ok := this.regMap[no]

	return c, ok
}

func (this *ServiceManager) Drop(no int64) {
	this.mapLock.Lock()
	defer this.mapLock.Unlock()
	delete(this.regMap, no)
}

func sendFPkg(c io.Writer, msg *smn_base.FPkg) error {
	bts, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = smn_net.WriteBytes(bts, c)

	return err
}

func (this *ServiceManager) Work() {
	go func() {
	nodeFor:
		for {
			select {
			case <-this.close:
				break nodeFor
			case msg := <-this.sendChan:
				err := sendFPkg(this.conn, msg)
				if err != nil {
					this.OnErr(err)
				}
			}
		}

		_ = this.conn.Close()
	}()
	go func() {
		for {
			bts, err := smn_net.ReadBytes(this.conn)
			if err != nil {
				this.OnErr(err)
			}

			pkg := &smn_base.FPkg{}
			err = proto.Unmarshal(bts, pkg)

			if err != nil {
				this.OnErr(err)
			}

			this.recv(pkg)
		}
	}()
}

func (this *ServiceManager) Close() {
	this.close <- 1
}

type ForwardConnItf interface {
	net.Conn
	SendToSM(msg []byte)
	RecvFromSM(msg []byte)
	Desc() string
	SetTimeOut(t time.Duration)
	ErrClose(err error)
}

type FConn struct {
	localAddr  net.Addr
	remoteAddr net.Addr
	no         int64
	mgr        *ServiceManager
	desc       string
	cache      *smn_stream.ByteCache
	err        error
}

func (this *FConn) Read(b []byte) (n int, err error) {
	return this.cache.Read(b)
}

func (this *FConn) Write(b []byte) (n int, err error) {
	if this.err != nil {
		return 0, this.err
	}

	this.SendToSM(b)

	return 0, nil
}

func (this *FConn) ErrClose(err error) {
	this.err = err
	this.mgr.Drop(this.no)
	this.cache.ErrorClose(err)
}
func (this *FConn) Close() error {
	this.mgr.Drop(this.no)
	this.cache.Close()

	return nil
}

func (this *FConn) LocalAddr() net.Addr {
	return this.localAddr
}

func (this *FConn) RemoteAddr() net.Addr {
	return this.remoteAddr
}

func (this *FConn) SetDeadline(t time.Time) error {
	return nil
}

func (this *FConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (this *FConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (this *FConn) SendToSM(msg []byte) {
	this.mgr.Send(&smn_base.FPkg{NO: this.no, Msg: msg})
}

func (this *FConn) RecvFromSM(msg []byte) {
	this.cache.Write(msg)
}

func (this *FConn) Desc() string {
	return this.desc
}

func (this *FConn) SetTimeOut(t time.Duration) {
	this.cache.SetTimeOut(t)
}
