package main

import (
	"flag"

	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/proto_compile"
)

var (
	comp       string
	protocPath string
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	o := flag.String("o", "./pb/", "output path")
	i := flag.String("i", "./datas/proto/", "input dir path.")
	goMod := flag.String("gm", "github.com/ProtossGenius/SureMoonNet", "go moudle.")
	lang := flag.String("lang", "go", "output language, cpp/csharp/java/javanano/objc/python/ruby")
	flag.Parse()
	checkerr(proto_compile.Compile(*i, *o, *goMod, *lang))
}
