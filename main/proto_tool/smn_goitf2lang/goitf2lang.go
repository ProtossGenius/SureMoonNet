package main

import (
	"flag"
	"os"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/gitf2lang"
)

var GoItfToLang = map[string]gitf2lang.FuncGoItfToLang{
	"cpp":  gitf2lang.WriteCppPkg,
	"java": gitf2lang.WriteJavaPkg,
}

func writeInterface(lang, out, pkg string, list []*smn_pglang.ItfDef) {
	f, ok := GoItfToLang[lang]
	if !ok {

	}
	f(out, pkg, list)
}

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
	for pkg, list := range itfs {
		writeInterface(*lang, *o, pkg, list)
	}
}
