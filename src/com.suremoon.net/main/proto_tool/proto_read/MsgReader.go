package main

import (
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/smn/analysis/proto_msg_map"
	"com.suremoon.net/smn/code_file_build"
	"flag"
	"os"
	"strings"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	protoPath := flag.String("proto", "./datas/proto/", "proto file's path")
	pkgHead := flag.String("pkgh", "pb/", "proto's pkg head")
	o := flag.String("o", "./src/pbr/read.go", "out path.")
	flag.Parse()
	err := os.MkdirAll((*o)[:strings.LastIndex(*o, "/")], os.ModePerm)
	checkerr(err)
	file, err := smn_file.CreateNewFile(*o)
	checkerr(err)
	list, cnm, err := proto_msg_map.Dict(*protoPath)
	checkerr(err)
	fileWriter := code_file_build.NewGoFile("smn_pbr", file, "product by tools, should not change this file.", "Author: SureMoon", "")
	fileWriter.AddImports(code_file_build.LocalImportable("./src"))
	fileWriter.Import("github.com/golang/protobuf/proto")
	fileWriter.Import("pb/dict")
	funcList := fileWriter.AddBlock("var funcList = []funcGetMsg")
	for _, pm := range list {
		constName := pm.Name
		if constName == "None" || strings.HasPrefix(constName, "//") {
			continue
		}
		funcList.WriteLine("dict.EDict_%s:%s,", pm.Name, pm.Name)
		f := fileWriter.AddBlock("func %s(bytes []byte) proto.Message {", constName)
		clzName := cnm[constName]
		f.Imports(*pkgHead + strings.Split(clzName, ".")[0])
		f.WriteLine("msg := &%s{}", clzName)
		f.WriteLine("proto.Unmarshal(bytes, msg)")
		f.WriteLine("return msg")
	}
	fileWriter.WriteLine("type funcGetMsg func(bytes []byte) proto.Message")

	fileWriter.WriteLine(`func GetMsgByDict(bytes []byte, dict dict.EDict) proto.Message {
	dictId := int(dict)
	if dictId >= len(funcList) || dictId < 0 || funcList[dictId] == nil {
		return nil
	}
	return funcList[dictId](bytes)
}

`)

	fileWriter.Output()
}
