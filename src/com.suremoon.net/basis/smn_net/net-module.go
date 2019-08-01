package smn_net

import (
	"fmt"
	"github.com/golang/protobuf/proto"
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

type MessageAdapterItf interface {
	WriteMessage(message proto.Message) error
	ReadMessage() ([]byte, error)
}

type MessageAdapter struct {
	c net.Conn
}

func NewMessageAdapter(conn net.Conn) MessageAdapterItf {
	return &MessageAdapter{c: conn}
}

func (this *MessageAdapter) WriteMessage(message proto.Message) error {
	bytes, err := proto.Marshal(message)
	if iserr(err) {
		return err
	}
	err = WriteInt(len(bytes), this.c)
	if iserr(err) {
		return err
	}
	_, err = this.c.Write(bytes)
	return err
}

func (this *MessageAdapter) ReadMessage() ([]byte, error) {
	len, err := ReadInt(this.c)
	if iserr(err) {
		return nil, err
	}
	bytes := make([]byte, len)
	rl, err := this.c.Read(bytes)
	if err != nil {
		return nil, err
	}
	if rl != len {
		return nil, fmt.Errorf(ErrNotGetEnoughLengthBytes, len, rl)
	}
	return bytes, err
}
