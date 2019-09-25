package smn_net

import (
	"fmt"
	"net"
)

type TcpServer struct {
	Port       int
	Listener   net.Listener
	AcceptChan chan net.Conn
	ErrChan    chan error
	OnRunning  bool
	Data       interface{}
}

func NewTcpServer(port int, acceptSize int) (this *TcpServer, err error) {
	this = &TcpServer{Port: port, ErrChan: make(chan error, 50), AcceptChan: make(chan net.Conn, acceptSize)}
	this.Listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	return
}

func (this *TcpServer) Close() {
	this.OnRunning = false
}

func (this *TcpServer) Run() {
	this.OnRunning = true
	for this.OnRunning {
		conn, err := this.Listener.Accept()
		if err != nil {
			this.ErrChan <- err
		}
		this.AcceptChan <- conn
	}
	this.Listener.Close()
}
