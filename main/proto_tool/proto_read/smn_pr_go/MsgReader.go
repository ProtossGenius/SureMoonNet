package main

import (
	"flag"

	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/proto_read_lang"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	protoPath := flag.String("proto", "./datas/proto/", "proto file's path")
	module := flag.String("module", "github.com/ProtossGenius/SureMoonNet", "go mod")
	lang := flag.String("lang", "go", "output coding language.")
	flag.Parse()

	err := proto_read_lang.Write(*lang, *protoPath, *module)
	checkerr(err)
}
