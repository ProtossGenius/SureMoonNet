package proto_read_lang

import (
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/proto_msg_map"
	"github.com/ProtossGenius/SureMoonNet/smn/code_file_build"
)

func GoMsgReader(protoPath, pkgHead, goPath, ext string) (err error) {
	o := "./pbr/read.go"
	err = os.MkdirAll((o)[:strings.LastIndex(o, "/")], os.ModePerm)
	if err != nil {
		return err
	}
	file, err := smn_file.CreateNewFile(o)
	if err != nil {
		return err
	}
	list, cnm, err := proto_msg_map.Dict(protoPath)
	if err != nil {
		return err
	}
	fileWriter := code_file_build.NewGoFile("smn_pbr", file, "product by tools, should not change this file.", "Author: SureMoon", "")
	fileWriter.AddImports(code_file_build.LocalImptTarget(goPath, goPath+ext))
	fileWriter.Import("github.com/golang/protobuf/proto")
	fileWriter.Import("pb/smn_dict")
	funcList := fileWriter.AddBlock("var funcList = []funcGetMsg")
	for _, pm := range list {
		constName := pm.Name
		if constName == "None" || strings.HasPrefix(constName, "//") {
			continue
		}
		funcList.WriteLine("smn_dict.EDict_%s:%s,", pm.Name, pm.Name)
		f := fileWriter.AddBlock("func %s(bytes []byte) proto.Message {", constName)
		clzName := cnm[constName]
		f.Imports(pkgHead + strings.Split(clzName, ".")[0])
		f.WriteLine("msg := &%s{}", clzName)
		f.WriteLine("proto.Unmarshal(bytes, msg)")
		f.WriteLine("return msg")
	}
	fileWriter.WriteLine("type funcGetMsg func(bytes []byte) proto.Message")

	fileWriter.WriteLine(`func GetMsgByDict(bytes []byte, dict smn_dict.EDict) proto.Message {
	dictId := int(dict)
	if dictId >= len(funcList) || dictId < 0 || funcList[dictId] == nil {
		return nil
	}
	return funcList[dictId](bytes)
}

`)

	_, err = fileWriter.Output()
	return
}
