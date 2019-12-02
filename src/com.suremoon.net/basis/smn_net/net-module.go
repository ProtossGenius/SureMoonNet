package smn_net

import (
	"fmt"
	"net"
	"com.suremoon.net/basis/smn_err"
)

type OnAccept func(conn net.Conn)

func DftAccept(conn net.Conn) {
	panic("You Should Set AcceptFunc For smn_net.TcpServer")
}

type TcpServer struct {
	Port      int
	Listener  net.Listener
	OnAccept  OnAccept
	OnErr     smn_err.OnErr
	OnRunning bool
	Data      interface{}
}

func NewTcpServer(port int, acceptSize int) (this *TcpServer, err error) {
	this = &TcpServer{Port: port, OnErr: smn_err.DftOnErr, OnAccept: DftAccept}
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
			this.OnErr(err)
		}
		this.OnAccept(conn)
	}
	this.Listener.Close()
}
