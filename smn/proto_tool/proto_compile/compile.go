package proto_compile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/smn/analysis/proto_msg_map"
)

func protoHead(pkg string) string {
	return fmt.Sprintf("syntax = \"proto3\";\noption java_package = \"pb\";\noption java_outer_classname=\"%s\";\npackage %s;\n\n", pkg, pkg)

}

//生成字典协议
func Dict(in string) error {
	list, _, err := proto_msg_map.Dict(in)
	file, err := smn_file.CreateNewFile(in + "/smn_dict.proto")
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(protoHead("smn_dict"))
	file.WriteString("enum EDict{\n")
	for _, val := range list {
		file.WriteString(fmt.Sprintf("\t%s = %d;\n", val.Name, val.Id))
	}
	file.WriteString("}\n")
	return nil
}

func getPkg(path string) string {
	data, err := smn_file.FileReadAll(path)
	if err != nil {
		return ""
	}
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

type compileFunc func(in, out, exportPath, ignoreDir, comp string) error

var CompileMap = map[string]compileFunc{
	"cpp":  CppCompile,
	"java": CppCompile,
}

func DefautCompile(in, out, goMoudle, ignoreDir, comp string) error {
	var ret_err error
	smn_file.DeepTraversalDir(in, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".proto") {
			op := in + "/temp/" + getPkg(path)
			op2 := in + "/temp/" + goMoudle + getPkg(path)
			os.MkdirAll(op, os.ModePerm)
			os.MkdirAll(op2, os.ModePerm)
			data, err := smn_file.FileReadAll(path)
			if err != nil {
				ret_err = err
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			lines := strings.Split(string(data), "\n")
			file, err := smn_file.CreateNewFile(op + "/" + info.Name())
			if err != nil {
				ret_err = err
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			defer file.Close()
			file2, err := smn_file.CreateNewFile(op2 + "/" + info.Name())
			if err != nil {
				ret_err = err
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			for _, line := range lines {
				nl := strings.TrimSpace(line)
				if strings.HasPrefix(nl, "import") {
					nl = strings.Split(nl[6:], ";")[0]
					nl = strings.Replace(nl, "\"", "", -1)
					nl = strings.TrimSpace(nl)
					if smn_file.IsFileExist(in + "/" + nl) {
						line = strings.Replace(line, nl, goMoudle+getPkg(in+"/"+nl)+"/"+nl, -1)
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
	smn_file.DeepTraversalDir(in+"/temp/", func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() && info.Name() == ignoreDir {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		if strings.HasSuffix(info.Name(), ".proto") {
			var stderr bytes.Buffer
			c := exec.Command("protoc", fmt.Sprintf(comp, out), "-I", in+"/temp/", path)
			c.Stderr = &stderr
			err := c.Run()
			if err != nil {
				panic(fmt.Errorf("%s: %s", err.Error(), stderr.String()))
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	os.RemoveAll(in + "/temp")
	return ret_err
}

func CppCompile(in, out, goMoudle, ignoreDir, comp string) error {
	var ret_err error
	smn_file.DeepTraversalDir(in, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() && info.Name() == ignoreDir {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		if strings.HasSuffix(info.Name(), ".proto") {
			var stderr bytes.Buffer
			c := exec.Command("protoc", fmt.Sprintf(comp, out), "-I", in, path)
			c.Stderr = &stderr
			err := c.Run()
			if err != nil {
				ret_err = err
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return ret_err
}

func Compile(protoDir, codeOutPath, goMod, lang string) error {
	if !smn_file.IsFileExist(codeOutPath) {
		err := os.MkdirAll(codeOutPath, os.ModePerm)
		if err != nil {
			return err
		}
	}
	Dict(protoDir)
	comp := "--" + lang + "_out=%s" //"--go_out=%s"
	extPath := strings.Replace(goMod, "\\", "/", -1) + "/" + strings.Replace(codeOutPath, "./", "", -1)
	ignoreDir := strings.Split(extPath, "/")[0]
	if f, ok := CompileMap[lang]; ok {
		return f(protoDir, codeOutPath, extPath, ignoreDir, comp)
	} else {
		return DefautCompile(protoDir, codeOutPath, extPath, ignoreDir, comp)
	}
}
