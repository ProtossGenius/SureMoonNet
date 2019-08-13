package main

import (
	"flag"
	"fmt"
	"os"

	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/smn/analysis/smn_rpc_itf"
	"com.suremoon.net/smn/code_file_build"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/** file as:
package xxxx
import(...)

*/
func writeSvrRpcFile(path string, list []*smn_pglang.ItfDef) {
	code_file_build.NewGoFile("", nil)
	for _, itf := range list {
		fmt.Println(itf.Package, "============")
	}
}

func writeClientRpcFile(path string, list []*smn_pglang.ItfDef) {

}

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc insterface;'s net accepter, from proto.Message call interface.")
	s := flag.Bool("s", false, "is product server code")
	c := flag.Bool("c", false, "is product client code")
	flag.Parse()
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	check(err)
	for _, list := range itfs {
		if *s {
			op := *o + "/svrrpc/"
			err := os.MkdirAll(op, os.ModePerm)
			check(err)
			writeSvrRpcFile(op, list)
		}
		if *c {
			op := *o + "/clientrpc/"
			err := os.MkdirAll(op, os.ModePerm)
			check(err)
			writeClientRpcFile(op, list)
		}
	}
}
