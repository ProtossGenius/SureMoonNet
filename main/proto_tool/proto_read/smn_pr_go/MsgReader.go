package main

import (
	"flag"
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/proto_read_lang"
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
	goPath := flag.String("gopath", "$GOPATH", "go path")
	ext := flag.String("ext", "", "exturn path.")
	lang := flag.String("lang", "go", "output coding language.")
	flag.Parse()
	f, ok := ReaderMap[*lang]
	if !ok {
		panic(fmt.Errorf("Error! not support language <%s>, you can contact us to achieve. ", *lang))
	}
	err := f(*protoPath, *pkgHead, *goPath, *ext)
	checkerr(err)
}
