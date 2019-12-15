package port_forward

import (
	"net"
	"errors"
	"io"
	"log"
)

const (
	ERR_IN_SOCKET_IS_NIL ="ERR_IN_SOCKET_IS_NIL"
	ERR_OUT_SOCKET_IS_NIL =  "ERR_OUT_SOCKET_IS_NIL"
)

type PortForwardWorker struct {
	in net.Conn
	out net.Conn
	pc chan int
}

func NewPortForwardWorker() *PortForwardWorker{
	return &PortForwardWorker{pc:make(chan int, 2)}
}

func (this *PortForwardWorker) DoWork() error{
	if this.in == nil {
		return errors.New(ERR_IN_SOCKET_IS_NIL)
	}
	if this.out == nil{
		return errors.New(ERR_OUT_SOCKET_IS_NIL)
	}
	go func() {
		this.pc<-1
		_, err := io.Copy(this.in, this.out)
		if err != nil{
			log.Println(err)
		}
	}()
	go func() {
		this.pc<-1
		_, err := io.Copy(this.out, this.in)
		if err != nil{
			log.Println(err)
		}
	}()
	return nil
}

func (this *PortForwardWorker) Wait() {
	<-this.pc
	<-this.pc
}

func (this *PortForwardWorker) SetInOut(in, out net.Conn){
	this.SetOut(out)
	this.SetIn(in)
}

func (this *PortForwardWorker) SetIn(s net.Conn){
	if s != nil{
		this.in = s
	}
}
func (this *PortForwardWorker) SetOut(s net.Conn){
	if s != nil{
		this.out = s
	}
}