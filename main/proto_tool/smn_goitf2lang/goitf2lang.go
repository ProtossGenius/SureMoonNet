package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/smn_rpc_itf"
)

var CodePkgConst = map[string]string{
	"cpp": "#include <vector>\n\nnamespace %s{",
}

var CodeEndConst = map[string]string{
	"cpp": "}",
}

type ItfTrans func(itf *smn_pglang.ItfDef) string

func Go2Cpp(itf *smn_pglang.ItfDef) string {
	lines := []string{}
	add := func(str ...string) {
		lines = append(lines, str...)
	}
	getReturns := func(itfName string, rets []*smn_pglang.VarDef) string {
		switch len(rets) {
		case 0:
			return "void"
		case 1:
			if rets[0].ArrSize != 0 {
				return fmt.Sprintf("std::Vector<%s>", rets[0].Type)
			}
			return rets[0].Type
		default:
			return "unknow"
		}
	}
	indentation := strings.Repeat("\t", 1)
	//interface(class) name
	add(fmt.Sprintf("%sclass %s {", indentation, itf.Name))
	for _, f := range itf.Functions {
		indent2 := indentation + "\t"
		add(fmt.Sprintf(indent2+"virtual %s %s() = 0;", getReturns(itf.Name, f.Returns), f.Name))
	}

	//}
	add(fmt.Sprintf("%s}", indentation))
	return strings.Join(lines, "\n")
}

var GoItfToLang = map[string]ItfTrans{
	"cpp": Go2Cpp,
}

func writeInterface(lang, out, pkg string, list []*smn_pglang.ItfDef) {
	converter, ok := GoItfToLang[lang]
	if !ok {
		panic(fmt.Errorf(`language %s now not support, you can goto github.com/ProtossGenius/SureMoonNet help write it.
		go file path is SureMoonNet/main/proto_tool/smn_goitf2lang/goitf2lang.go`, lang))
	}
	f, err := smn_file.CreateNewFile(out + "/" + pkg + ".h")
	checkerr(err)
	defer f.Close()
	write := func(str string) {
		_, err := f.WriteString(str)
		checkerr(err)
	}
	writeln := func(str string) {
		write(str + "\n")
	}
	writef := func(str string, a ...interface{}) {
		write(fmt.Sprintf(str, a...))
	}
	//write package(in cpp is namespace.)
	if f, ok := CodePkgConst[lang]; ok {
		writef(f+"\n", pkg)
	}
	for _, itf := range list {
		writeln(converter(itf))
	}
	//write file end.
	if f, ok := CodeEndConst[lang]; ok {
		writeln(f)
	}
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
