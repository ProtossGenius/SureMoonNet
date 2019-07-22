package smn_net

import (
	"fmt"
	"net"
)

type SSocketAcceptFunc func(svr *ServerSocket, conn net.Conn)

type ServerSocket struct {
	AcceptFunc SSocketAcceptFunc
	Port       int
	Listener   net.Listener
	ErrChan    chan error
	Data       interface{}
}

func NewServer(port int, acceptFunc SSocketAcceptFunc) (this *ServerSocket, err error) {
	this = &ServerSocket{Port: port, AcceptFunc: acceptFunc, ErrChan: make(chan error, 50)}
	this.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	return
}

func (this *ServerSocket) Run() {
	for {
		conn, err := this.Listener.Accept()
		if err != nil {
			this.ErrChan <- err
		}
		go this.AcceptFunc(this, conn)
	}
	this.Listener.Close()
}
