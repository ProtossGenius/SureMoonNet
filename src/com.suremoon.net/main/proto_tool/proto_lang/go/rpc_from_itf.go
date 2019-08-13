package main

import (
	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/smn/analysis/smn_rpc_itf"
	"flag"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func writeRpcFile(path string, list []*smn_pglang.ItfDef) {
}

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc insterface;'s net accepter, from proto.Message call interface.")
	flag.Parse()
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	check(err)
	for pkg, list := range itfs {
		writeRpcFile(*o+"/rpc_"+pkg+".go", list)
	}
}
