package muti_service

import (
	"pb/smn_base"
	"net"
	"com.suremoon.net/basis/smn_net"
	"fmt"
	"github.com/golang/protobuf/proto"
	"com.suremoon.net/basis/smn_err"
	"time"
	"sync"
)

const (
	ErrServerNumNotExist = "ErrServerNumNotExist :[%d]"
)

type ServiceManager struct {
	OnErr    smn_err.OnErr
	TimeOut  time.Duration
	conn     net.Conn
	sendChan chan *smn_base.FPkg
	regMap   map[int64]ForwardConnItf //key always > 0
	close    chan int
	mapLock  sync.Mutex //因为map操作并不频繁，所以不用计较锁的开销
}

func NewServiceManager(conn net.Conn) *ServiceManager {
	return &ServiceManager{conn: conn, close: make(chan int, 1), regMap: make(map[int64]ForwardConnItf), sendChan: make(chan *smn_base.FPkg, 1024), OnErr: smn_err.DftOnErr, TimeOut: 1 * time.Second}
}

func (this *ServiceManager) send(p *smn_base.FPkg) {
	this.sendChan <- p
}

func (this *ServiceManager) recv(p *smn_base.FPkg) {
	conn, ok := this.GetFConn(p.NO)
	if ok {
		conn.recv(p.Msg)
	} else {
		this.OnErr(fmt.Errorf(ErrServerNumNotExist, p.NO))
	}
}

func (this *ServiceManager) Regitster(no int64, desc string) (conn ForwardConnItf, isExist bool) {
	this.mapLock.Lock()
	defer this.mapLock.Unlock()
	if c, ok := this.regMap[no]; ok {
		return c, true
	}
	conn = &fConn{no: no, mgr: this, desc: desc, conn: this.conn, recvChan: make(chan []byte, 1000), TimeOut: 1 * time.Second}
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

func sendFPkg(c net.Conn, msg *smn_base.FPkg) error {
	bts, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	err = smn_net.WriteInt(len(bts), c)
	if err != nil {
		return err
	}
	_, err = c.Write(bts)
	return err
}

func (this *ServiceManager) Work() {
	go func() {
		for {
			select {
			case <-this.close:
				break
			case msg := <-this.sendChan:
				err := sendFPkg(this.conn, msg)
				if err != nil {
					this.OnErr(err)
				}
			}
		}
		this.conn.Close()
	}()
	go func() {
		for {
			size, err := smn_net.ReadInt(this.conn)
			if err != nil {
				this.OnErr(err)
			}
			bts := make([]byte, size)
			rs, err := this.conn.Read(bts)
			if err != nil {
				this.OnErr(err)

			}
			pkg := &smn_base.FPkg{}
			if rs != size {
				this.OnErr(fmt.Errorf(smn_net.ErrNotGetEnoughLengthBytes, size, rs))
			}
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
	send(msg []byte)
	recv(msg []byte)
	Desc()string
}

type fConn struct {
	conn     net.Conn
	no       int64
	mgr      *ServiceManager
	desc     string
	recvChan chan []byte
	last     []byte
	readLock sync.Mutex
	TimeOut  time.Duration // timeout should not be zero
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

//copy `size` byte from fConn.last to b[start:]
func (this *fConn) readCopy(b []byte, start, size int) (copyLen int) {
	copyLen = minInt(size, len(this.last))
	pos := start
	for i := 0; i < copyLen; i++ {
		b[pos] = this.last[i]
		pos++
	}
	this.last = this.last[copyLen:]
	return
}

func (this *fConn) Read(b []byte) (n int, err error) {
	this.readLock.Lock()
	defer this.readLock.Unlock()
	size := len(b)
	acLen := 0
	if this.TimeOut == 0{
		this.TimeOut = 9 * time.Second
	}
	tout := time.Tick(this.TimeOut)
	for size > 0 {
		copyLen := this.readCopy(b, acLen, size)
		acLen += copyLen
		size -= copyLen
		if size == 0 {
			break
		}
		select {
		case this.last = <-this.recvChan:
		case <-tout:
			return acLen, nil
		}
	}
	return len(b), nil
}

func (this *fConn) Write(b []byte) (n int, err error) {
	this.send(b)
	return 0, nil
}

func (this *fConn) Close() error {
	this.mgr.Drop(this.no)
	return nil
}

func (this *fConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *fConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *fConn) SetDeadline(t time.Time) error {
	return nil
}

func (this *fConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (this *fConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (this *fConn) send(msg []byte) {
	this.mgr.send(&smn_base.FPkg{NO: this.no, Msg: msg})
}

func (this *fConn) recv(msg []byte) {
	this.recvChan<-msg
}

func (this *fConn) Desc() string {
	return this.desc
}