package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/itf2rpc"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc interface's net accepter, from proto.Message call interface.")
	s := flag.Bool("s", true, "is product server code")
	c := flag.Bool("c", true, "is product client code")
	pMod := flag.String("module", "github.com/ProtossGenius/SureMoonNet", "go module.")
	flag.Parse()

	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	check(err)

	for path, list := range itfs {
		if len(list) == 0 {
			continue
		}

		fullPath, err := filepath.Abs(path)
		check(err)
		pwdPath, err := filepath.Abs("./")
		check(err)
		//get fullPkg
		fullPkg := *pMod + strings.Replace(fullPath, pwdPath, "", -1)

		for _, itf := range list {
			pkg := itf.Package

			if *s {
				op := *o + "/svrrpc/svr_rpc_" + pkg + "/"
				err := os.MkdirAll(op, os.ModePerm)
				check(err)

				err = itf2rpc.Write("go_s", op, *pMod, fullPkg, itf)
				check(err)
			}

			if *c {
				op := *o + "/cltrpc/clt_rpc_" + pkg + "/"
				err := os.MkdirAll(op, os.ModePerm)
				check(err)

				err = itf2rpc.Write("go_c", op, *pMod, fullPkg, itf)
				check(err)
			}
		}
	}
}
