package main

import (
	"flag"
	"com.suremoon.net/smn/proto_tool/proto_read_lang"
	"fmt"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

var ReaderMap = map[string]proto_read_lang.MsgReader{
	"go": proto_read_lang.GoMsgReader,
}

func main() {
	protoPath := flag.String("proto", "./datas/proto/", "proto file's path")
	pkgHead := flag.String("pkgh", "pb/", "proto's pkg head")
	o := flag.String("o", "./src/pbr/read.go", "out path.")
	lang := flag.String("lang", "go", "output coding language.")
	flag.Parse()
	f, ok := ReaderMap[*lang]
	if !ok {
		panic(fmt.Errorf("Error! not support language <%s>, you can contact us to achieve. ", *lang))
	}
	err := f(*protoPath, *pkgHead, *o)
	checkerr(err)
}
