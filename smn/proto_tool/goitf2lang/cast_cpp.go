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
		if goType == "int" || goType == "uint" {
			return goType + "64" + "_t"
		}
		return goType + "_t"
	}
	switch goType {
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

func TooCppRet(rets []*smn_pglang.VarDef, pkg, itfName, fName string) string {
	if len(rets) == 0 {
		return "void"
	}
	if len(rets) == 1 {
		p := ToCppVarDef(rets[0])
		return p.Type
	}
	return fmt.Sprintf("rip_%s::%s_%s_Ret", pkg, itfName, fName)

}

func WriteCppItf(out, pkg string, itf *smn_pglang.ItfDef) {
	dir := out + "/smn_itf/"
	if !smn_file.IsFileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		checkerr(err)
	}
	dir += pkg + "."
	f, err := smn_file.CreateNewFile(dir + itf.Name + ".h")
	checkerr(err)
	defer f.Close()
	writef := func(s string, a ...interface{}) {
		_, err := f.WriteString(fmt.Sprintf(s, a...))
		checkerr(err)
	}
	writef(`#pragma once
#include<vector>
#include"smncpp/socket_itf.h"
`)

	for _, f := range itf.Functions {
		if len(f.Returns) > 1 {
			writef("#include \"pb/rip_%s.pb.h\"", pkg)
		}
	}

	writef(`
namespace %s{

`, pkg)
	defer writef("}//namespace %s", pkg)

	writef("class %s {\npublic:\n", itf.Name)
	defer writef("};\n")
	for _, f := range itf.Functions {
		writef("\tvirtual %s %s(%s) = 0;\n", TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name, ToCppParam(f.Params))
	}
}

func WriteCppPkg(out, pkg string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		WriteCppItf(out, pkg, itf)
	}
}
