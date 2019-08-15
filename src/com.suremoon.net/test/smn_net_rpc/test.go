package main

import (
	"com.suremoon.net/basis/smn_net"
	"fmt"
	"net"
	"rpc_nitf/clientrpc"
	"rpc_nitf/svrrpc"
	"time"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type login struct {
}

func (this *login) DoLogin(user, pswd string, code int) (bool, int) {
	fmt.Println(user, pswd, code)
	return false, 0
}

func (this *login) Test1(a []string, b []int, c []uint, d []uint64, e []int32) []int {
	panic("implement me")
}

func AccpterRun(adapter smn_net.MessageAdapterItf) {
	rpcSvr := svr_rpc_rpc_itf.NewSvrRpcLogin(&login{})
	for {
		msg, err := adapter.ReadCall()
		check(err)
		dict, res, err := rpcSvr.OnMessage(msg)
		adapter.WriteRet(dict, res, err)
	}
}

func accept(c chan net.Conn) {
	for {
		conn := <-c
		adapter := smn_net.NewMessageAdapter(conn)
		go AccpterRun(adapter)
	}
}

func RunSvr() {
	svr, err := smn_net.NewTcpServer(1000, 100)
	check(err)
	go accept(svr.AcceptChan)
	svr.Run()
}

func main() {
	go RunSvr()
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:1000")
	check(err)
	client := clt_rpc_rpc_itf.NewCltRpcLogin(smn_net.NewMessageAdapter(conn))
	b, i := client.DoLogin("user---", "pswd_____", -1)
	fmt.Println("cccccc   ", b, i)
	client.Test1(nil, nil, nil, nil, nil)
}
