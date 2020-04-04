package smn_port_forward

import (
	"errors"
	"io"
	"net"
)

const (
	ERR_IN_SOCKET_IS_NIL  = "ERR_IN_SOCKET_IS_NIL"
	ERR_OUT_SOCKET_IS_NIL = "ERR_OUT_SOCKET_IS_NIL"
)

type PortForwardWorker struct {
	in          net.Conn //client connect in
	out         net.Conn //connect to outside (remote)
	pc          chan int
	CopyFromIn  func(dst io.Writer, src io.Reader) (int64, error)
	CopyFromOut func(dst io.Writer, src io.Reader) (int64, error)
}

func NewPortForwardWorker() *PortForwardWorker {
	return &PortForwardWorker{pc: make(chan int, 2), CopyFromIn: io.Copy, CopyFromOut: io.Copy}
}

func (this *PortForwardWorker) DoWork(dealErr func(err error)) error {
	if this.in == nil {
		return errors.New(ERR_IN_SOCKET_IS_NIL)
	}
	if this.out == nil {
		return errors.New(ERR_OUT_SOCKET_IS_NIL)
	}
	go func() {
		this.pc <- 1
		defer func() {
			<-this.pc
			this.in.Close()
			this.out.Close()
		}()
		_, err := this.CopyFromOut(this.in, this.out)
		dealErr(err)
	}()
	go func() {
		this.pc <- 1
		defer func() { <-this.pc }()
		_, err := this.CopyFromIn(this.out, this.in)
		dealErr(err)
	}()
	return nil
}

func (this *PortForwardWorker) Wait() {
	<-this.pc
	<-this.pc
}

func (this *PortForwardWorker) SetInOut(in, out net.Conn) {
	this.SetOut(out)
	this.SetIn(in)
}

func (this *PortForwardWorker) SetIn(s net.Conn) {
	if s != nil {
		this.in = s
	}
}
func (this *PortForwardWorker) SetOut(s net.Conn) {
	if s != nil {
		this.out = s
	}
}
