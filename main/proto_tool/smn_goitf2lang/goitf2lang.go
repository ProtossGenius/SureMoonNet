package main

import (
	"flag"
	"os"

	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/goitf2lang"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	i := flag.String("i", "./rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./cpp_itf/", "rpc needs proto output.")
	lang := flag.String("lang", "cpp", "lang to translate, now support [cpp]")
	flag.Parse()
	err := os.MkdirAll(*o, os.ModePerm)
	checkerr(err)
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	checkerr(err)
	for _, list := range itfs {
		pkg := list[0].Package
		goitf2lang.WriteInterface(*lang, *o, pkg, list)
	}
}
