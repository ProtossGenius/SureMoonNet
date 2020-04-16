package main

import (
	"flag"
	"os"

	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/itf2proto"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	i := flag.String("i", "./rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./datas/proto/", "proto output dir.")
	flag.Parse()
	err := os.MkdirAll(*o, os.ModePerm)
	checkerr(err)
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	checkerr(err)
	for _, list := range itfs {
		itf2proto.WriteProto(*o, list)
	}
}
