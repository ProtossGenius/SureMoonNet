package main

import (
	"fmt"
	"io"
	"net"
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go"
	"time"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func Svr(laddr string) {
	lis, err := kcp.ListenWithOptions(laddr, nil, 10, 3)
	checkerr(err)
	for {
		conn, e := lis.AcceptKCP()
		checkerr(e)
		fmt.Println(conn.RemoteAddr())
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
	go Svr(":10000")
	clt, err := kcp.DialWithOptions("localhost:10000", nil, 10, 3)
	checkerr(err)
	clt.Write([]byte("hello!!!!!11111111111111111111111111"))
	clt.Write([]byte("hello!!!!!2222222222222222222222222222222"))
	for {
		time.Sleep(1 * time.Second)
	}
}
