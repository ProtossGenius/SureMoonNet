package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/ProtossGenius/SureMoonNet/rpc_nitf/svrrpc/svr_rpc_rpc_itf"
	"github.com/ProtossGenius/SureMoonNet/test/rpc_itfs/rpc_itf"
	"github.com/gogo/protobuf/proto"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_net"
	"github.com/ProtossGenius/SureMoonNet/smn/net_libs/smn_rpc"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type login struct {
	rpc_itf.Login
}

func (l *login) DoLogin(user, pswd string, code int) (bool, int) {
	fmt.Println(user, pswd, code)
	return false, 0
}

func (l *login) Test1(a []string, b []int, c []uint, d []uint64, e []int32) []int {
	fmt.Println("call login.Test1")
	return []int{1, 2, 3, 4, 54321}
}

func (l *login) Test2(key string, c net.Conn) bool {
	if key == "hello" {
		_, err := smn_net.WriteString("world", c)
		check(err)

		return true
	}

	_, err := smn_net.WriteString("where is hello?", c)
	check(err)

	return false
}

func AccpterRun(adapter smn_rpc.MessageAdapterItf) {
	/*	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}

		adapter.Close()
	}()*/

	rpcSvr := svr_rpc_rpc_itf.NewSvrRpcLogin(&login{})

	for {
		msg, err := adapter.ReadCall()
		fmt.Println("get Msg : ", msg)
		check(err)
		dict, res, err := rpcSvr.OnMessage(msg, adapter.GetConn())
		fmt.Println("on Message dict = ", dict, ", res = ", fmt.Sprintf("%v", res), "err = ", err)
		if res != nil {
			bts, err := proto.Marshal(res)
			fmt.Println("res.Length = ", len(bts), "[", bts, "]", err)
		}
		check(err)
		_, err = adapter.WriteRet(dict, res, err)
		check(err)
	}
}

func accept(conn net.Conn) {
	adapter := smn_rpc.NewMessageAdapter(conn)
	go AccpterRun(adapter)
}

func RunSvr(port int) {
	svr, err := smn_net.NewTcpServer(port, accept)
	check(err)
	svr.Run()
}

func main() {
	pPort := flag.Int("port", 7000, "port.")
	flag.Parse()

	fmt.Printf("server run on port %d\n", *pPort)
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	RunSvr(*pPort)
}
