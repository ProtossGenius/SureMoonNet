package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_muti_write_cache"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func writeProto(oPath string, list []*smn_pglang.ItfDef) {
	pkg := list[0].Package
	file, err := smn_file.CreateNewFile(oPath)
	checkerr(err)
	w := smn_muti_write_cache.NewFileMutiWriteCache()
	writeLine := func(f string, a ...interface{}) {
		w.WriteTail(fmt.Sprintf(f, a...) + "\n")
	}
	w.Append(smn_muti_write_cache.NewStrCache(fmt.Sprintf("syntax = \"proto3\";\n\npackage rip_%s;\n", pkg)))
	impts := smn_muti_write_cache.NewStrCache()
	impMap := make(map[string]bool)
	checkImport := func(typ string) {
		if typ == "net.Conn" {
			return
		}
		list := strings.Split(typ, ".")
		if len(list) == 1 {
			return
		}
		if !impMap[list[0]] {
			impts.WriteTail(fmt.Sprintf("import \"%s.proto\";", list[0]) + "\n")
			impMap[list[0]] = true
		}
	}
	w.Append(impts)
	for _, itf := range list {
		for _, mtd := range itf.Functions {
			if mtd.Name[0] < 'A' || mtd.Name[0] > 'Z' {
				log.Printf("warning! mtd.Name %s's first letter not upper\n", mtd.Name)
				continue
			}
			writeLine("message %s_%s_Prm {", itf.Name, mtd.Name)
			for i, prm := range mtd.Params {
				isArray, typ := smn_str.ProtoUseDeal(prm.Type)
				if typ == "net.Conn" {
					continue
				}
				checkImport(typ)
				rpt := ""
				if isArray {
					rpt = "repeated "
				}
				writeLine("\t%s%s %s = %d;", rpt, typ, prm.Var, i+1)
			}
			writeLine("}")
			writeLine("message %s_%s_Ret {", itf.Name, mtd.Name)
			for i, res := range mtd.Returns {
				isArray, typ := smn_str.ProtoUseDeal(res.Type)
				checkImport(typ)
				rpt := ""
				if isArray {
					rpt = "repeated "
				}
				writeLine("\t%s%s %s = %d;", rpt, typ, res.Var, i+1)
			}
			writeLine("}")
		}
	}
	w.Output(file)
}

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./datas/proto/", "rpc needs proto output.")
	flag.Parse()
	err := os.MkdirAll(*o, os.ModePerm)
	checkerr(err)
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	checkerr(err)
	for pkg, list := range itfs {
		writeProto(*o+"/rip_"+pkg+".proto", list)
	}
}
