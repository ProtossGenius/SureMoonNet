package smn_net

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"pb/base"
	"pb/dict"
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
	WriteMessage(dict dict.EDict, message proto.Message) (int, error)
	ReadMessage() (*base.Call, error)
}

type MessageAdapter struct {
	c net.Conn
}

func NewMessageAdapter(conn net.Conn) MessageAdapterItf {
	return &MessageAdapter{c: conn}
}

func (this *MessageAdapter) WriteMessage(dict dict.EDict, message proto.Message) (int, error) {
	bytes, err := proto.Marshal(message)
	if iserr(err) {
		return 0, err
	}
	msg := &base.Call{Dict: dict, Msg: bytes}
	bytes, err = proto.Marshal(msg)
	err = WriteInt(len(bytes), this.c)
	if iserr(err) {
		return 0, err
	}
	return this.c.Write(bytes)
}

func (this *MessageAdapter) ReadMessage() (*base.Call, error) {
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
	msg := &base.Call{}
	proto.Unmarshal(bytes, msg)
	return msg, err
}
