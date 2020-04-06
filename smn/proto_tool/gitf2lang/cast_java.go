package gitf2lang

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
)

func ToJavaType(goType string) string {
	switch goType {
	case "int32", "uint32":
		return "int"
	case "int", "uint", "int64", "uint64":
		return "long"
	case "float32":
		return "float"
	case "float64":
		return "double"
	case "string":
		return "String"
	case "bool":
		return "boolean"
	case "net.Conn":
		return "com.suremoon.net.ConnFunc"
	}
	if strings.Contains(goType, "*") {
		goType = strings.Replace(goType, "*", "", -1)
		goType = "pb." + goType
	}
	return goType
}

func JavaBuiltInType(t string) bool {
	switch t {
	case "int", "short", "long", "String", "float", "double", "char", "byte", "boolean",
		"Interger", "Short", "Long", "Float", "Double", "Charecter", "Byte", "Boolean":
		return true
	}
	return false
}

func ToJavaVarDef(vd *smn_pglang.VarDef) *smn_pglang.VarDef {
	res := &smn_pglang.VarDef{Var: vd.Var, Type: vd.Type, ArrSize: vd.ArrSize}
	cnt := strings.Count(res.Type, "[]")
	res.Type = strings.Replace(res.Type, "[]", "", -1)
	res.Type = ToJavaType(res.Type)
	res.Type = res.Type + strings.Repeat("[]", cnt)
	return res
}
func ToJavaRet(param []*smn_pglang.VarDef, pkg, itfName, fName string) string {
	if len(param) == 0 {
		return "void"
	}
	if len(param) == 1 {
		vd := ToJavaVarDef(param[0])
		return vd.Type
	}
	return fmt.Sprintf("pb.%s.%s_%s_Ret", pkg, itfName, fName)
}
func ToJavaParam(param []*smn_pglang.VarDef) string {
	list := make([]string, len(param))
	for i, p := range param {
		p = ToJavaVarDef(p)
		if p.Var == "" {
			p.Var = fmt.Sprintf("sm_p%d", i)
		}
		list[i] = fmt.Sprintf("%s %s", p.Type, p.Var)
	}
	return strings.Join(list, ", ")

}

func WriteJavaItf(out, pkg string, itf *smn_pglang.ItfDef) {
	dir := out + "/smn_itf/" + pkg + "/"
	if !smn_file.IsFileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		checkerr(err)

	}
	f, err := smn_file.CreateNewFile(dir + itf.Name + ".java")
	checkerr(err)
	defer f.Close()
	writef := func(s string, a ...interface{}) {
		_, err := f.WriteString(fmt.Sprintf(s, a...))
		checkerr(err)
	}
	writef("package smn_itf.%s;\n//product by auto-code tools, you should never change it .\n//author SureMoonNet\n\n", pkg)

	writef("public interface %s {\n", itf.Name)
	defer writef("}\n")
	for _, f := range itf.Functions {
		writef("\t%s %s(%s);\n", ToJavaRet(f.Returns, pkg, itf.Name, f.Name), f.Name, ToJavaParam(f.Params))
	}
}

func WriteJavaPkg(out, pkg string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		WriteJavaItf(out, pkg, itf)
	}
}
