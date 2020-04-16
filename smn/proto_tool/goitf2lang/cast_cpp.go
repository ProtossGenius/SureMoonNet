package goitf2lang

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
)

func ToCppType(goType string) string {
	if strings.HasPrefix(goType, "int") || strings.HasPrefix(goType, "uint") {
		return goType + "_t"
	}
	switch goType {
	case "int":
		return "int64_t"
	case "uint":
		return "uint64_t"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "std::string"
	case "net.Conn":
		return "smnet::Conn"
	}
	if strings.Contains(goType, "*") {
		goType = strings.Replace(goType, "*", "", -1)
		goType = strings.Replace(goType, ".", "::", -1)
	}
	return goType
}

func CppBuiltInType(t string) bool {
	switch t {
	case "int", "unsigned int", "int32_t", "uint32_t", "long", "unsigned long", "long long", "unsigned long long",
		"int8_t", "uint8_t", "int16_t", "uint16_t", "int64_t", "uint64_t", "double", "float", "char", "unsigned char",
		"short", "unsigned short", "std::size_t":
		return true
	}
	return false
}

func ToCppVarDef(vd *smn_pglang.VarDef) *smn_pglang.VarDef {
	res := &smn_pglang.VarDef{Var: vd.Var}
	if vd.ArrSize != 0 {
		arrV := 0
		for strings.Contains(vd.Type, "[]") {
			vd.Type = strings.Replace(vd.Type, "[]", "", 1)
			arrV++
		}
		res.Type = fmt.Sprintf("%s%s%s", strings.Repeat("std::vector<", arrV), ToCppType(vd.Type), strings.Repeat(">", arrV))
	} else {
		res.Type = ToCppType(vd.Type)
	}
	return res
}

func ToCppParam(param []*smn_pglang.VarDef) string {
	if len(param) == 0 {
		return "void"
	}
	list := make([]string, len(param))
	for i, p := range param {
		p = ToCppVarDef(p)
		if p.Var == "" {
			p.Var = fmt.Sprintf("sm_p%d", i)
		}
		if CppBuiltInType(p.Type) {
			list[i] = fmt.Sprintf("%s %s", p.Type, p.Var)
		} else {
			list[i] = fmt.Sprintf("const %s& %s", p.Type, p.Var)
		}
	}
	return strings.Join(list, ", ")

}

func TooCppRet(param []*smn_pglang.VarDef, pkg, itfName, fName string) string {
	if len(param) == 0 {
		return "void"
	}
	if len(param) == 1 {
		return ToCppParam(param)
	}
	return fmt.Sprintf("%s::%s_%s_Ret", pkg, itfName, fName)

}

func WriteCppItf(out, pkg string, itf *smn_pglang.ItfDef) {
	dir := out + "/smn_itf/" + pkg + "/"
	if !smn_file.IsFileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		checkerr(err)
	}
	f, err := smn_file.CreateNewFile(dir + itf.Name + ".h")
	checkerr(err)
	defer f.Close()
	writef := func(s string, a ...interface{}) {
		_, err := f.WriteString(fmt.Sprintf(s, a...))
		checkerr(err)
	}
	writef(`#param once
#include<vector>
namespace %s{

`, pkg)
	defer writef("}//namespace %s", pkg)

	writef("class %s {\npublic:\n", itf.Name)
	defer writef("}\n")
	for _, f := range itf.Functions {
		writef("\tvirtual %s %s(%s) = 0;\n", TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name, ToCppParam(f.Params))
	}
}

func WriteCppPkg(out, pkg string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		WriteCppItf(out, pkg, itf)
	}
}
