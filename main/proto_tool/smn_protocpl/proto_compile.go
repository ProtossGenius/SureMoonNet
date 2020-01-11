package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/proto_msg_map"
)

var comp string

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func dict(in string) {
	list, _, err := proto_msg_map.Dict(in)
	file, err := smn_file.CreateNewFile(in + "smn_dict.proto")
	checkerr(err)
	file.WriteString("syntax = \"proto3\";\n\npackage smn_dict;\n\nenum EDict{\n")
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
	o := flag.String("o", "./pb/", "output path")
	i := flag.String("i", "./datas/proto/", "input dir path.")
	ep := flag.String("ep", "", "export path.")
	lang := flag.String("lang", "go", "output language, cpp/csharp/java/javanano/objc/python/ruby")
	flag.Parse()
	extPath := strings.Replace(*ep, "\\", "/", -1) + "/" + strings.Replace(*o, "./", "", -1)
	ignoreDir := strings.Split(extPath, "/")[0]
	err := os.MkdirAll(*o, os.ModePerm)
	checkerr(err)
	comp = "--" + *lang + "_out=%s" //"--go_out=%s"
	dict(*i)
	smn_file.DeepTraversalDir(*i, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			op := *i + "/temp/" + getPkg(path)
			op2 := *i + "/temp/" + extPath + getPkg(path)
			os.MkdirAll(op, os.ModePerm)
			os.MkdirAll(op2, os.ModePerm)
			data, err := smn_file.FileReadAll(path)
			checkerr(err)
			lines := strings.Split(string(data), "\n")
			file, err := smn_file.CreateNewFile(op + "/" + info.Name())
			file2, err := smn_file.CreateNewFile(op2 + "/" + info.Name())
			checkerr(err)
			for _, line := range lines {
				nl := strings.TrimSpace(line)
				if strings.HasPrefix(nl, "import") {
					nl = strings.Split(nl[6:], ";")[0]
					nl = strings.Replace(nl, "\"", "", -1)
					nl = strings.TrimSpace(nl)
					if smn_file.IsFileExist(*i + "/" + nl) {
						line = strings.Replace(line, nl, extPath+getPkg(*i+"/"+nl)+"/"+nl, -1)
					}
				}
				file.WriteString(line + "\n")
				file2.WriteString(line + "\n")
			}
			file.Close()
			file2.Close()
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	smn_file.DeepTraversalDir(*i+"/temp/", func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() && info.Name() == ignoreDir {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		if strings.HasSuffix(info.Name(), ".proto") {
			var stderr bytes.Buffer
			c := exec.Command(*p, fmt.Sprintf(comp, *o), "-I", *i+"/temp/", path)
			c.Stderr = &stderr
			//c := exec.Command(*p, fmt.Sprintf(comp, *o, *i, path))
			err := c.Run()
			if err != nil {
				panic(fmt.Errorf("%s: %s", err.Error(), stderr.String()))
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	os.RemoveAll(*i + "/temp")
}