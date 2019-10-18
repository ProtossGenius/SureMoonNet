package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func Svr() {
	lis, err := kcp.ListenWithOptions(":10000", nil, 10, 3)
	checkerr(err)
	for {
		conn, e := lis.AcceptKCP()
		checkerr(e)
		go func(conn net.Conn) {
			var buff = make([]byte, 1024, 1024)
			for {
				n, e := conn.Read(buff)
				if e != nil {
					if e == io.EOF {
						break
					}
					fmt.Println(errors.Wrap(e, "hello?"))
					break
				}
				fmt.Println("recv from client:", buff[:n])
			}
		}(conn)
	}
}

func main() {
	go Svr()
	clt, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	checkerr(err)
	clt.Write([]byte("hello!!!!!11111111111111111111111111"))
	clt.Write([]byte("hello!!!!!2222222222222222222222222222222"))
	fmt.Println(clt.RemoteAddr())
	time.Sleep(10 * time.Second)
}
