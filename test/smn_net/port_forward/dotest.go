package main

import (
	"net"
	"github.com/ProtossGenius/SureMoonNet/smn/net_libs/port_forward"
	"fmt"
)

func check(err error) {
	if err != nil{
		panic(err)
	}
}

func main() {
	svr, _ := net.Listen("tcp", ":10001")
	for{
		in, err := svr.Accept()
		check(err)
		fmt.Println("1111111111111111111111")
		out, err := net.Dial("tcp", "www.baidu.com:443")
		check(err)
		worker := port_forward.NewPortForwardWorker()
		worker.SetInOut(in, out)
		worker.DoWork()
	}
}
