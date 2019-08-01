package main

import (
	"bytes"
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/smn/analysis/proto_msg_map"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const comp = "--go_out=%s"

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func dict(in string) {
	list, _, err := proto_msg_map.Dict(in)
	file, err := smn_file.CreateNewFile(in + "dict.proto")
	checkerr(err)
	file.WriteString("syntax = \"proto3\";\n\npackage dict;\n\nenum EDict{\n")
	for _, val := range list {
		file.WriteString(fmt.Sprintf("\t%s = %d;\n", val.Name, val.Id))
	}
	file.WriteString("}\n")
	file.Close()
}

func getPkg(path string) string {
	data, err := smn_file.FileReadAll(path)
	checkerr(err)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {

		if strings.HasPrefix(line, "package") {
			pkg := strings.Split(line[7:], ";")[0]
			pkg = strings.TrimSpace(pkg)
			return pkg
		}
	}
	return ""
}

func main() {
	p := flag.String("p", "protoc", "protoc's path")
	o := flag.String("o", "./src/pb/", "output path")
	i := flag.String("i", "./datas/proto/", "input dir path.")
	flag.Parse()
	dict(*i)
	smn_file.DeepTraversalDir(*i, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			var stderr bytes.Buffer
			oPath := *o + getPkg(path) + "/"
			os.MkdirAll(oPath, os.ModePerm)
			c := exec.Command(*p, fmt.Sprintf(comp, oPath), "-I", *i, path)
			c.Stderr = &stderr
			//c := exec.Command(*p, fmt.Sprintf(comp, *o, *i, path))
			err := c.Run()
			if err != nil {
				panic(fmt.Errorf("%s: %s", err.Error(), stderr.String()))
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
}
