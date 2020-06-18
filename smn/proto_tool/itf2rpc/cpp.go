package itf2rpc

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/goitf2lang"
)

func ToCppType(goType string) string {
	return goitf2lang.ToCppType(goType)
}

func CppBuiltInType(t string) bool {
	return goitf2lang.CppBuiltInType(t)
}

func ToCppVarDef(vd *smn_pglang.VarDef) *smn_pglang.VarDef {
	return goitf2lang.ToCppVarDef(vd)
}

func ToCppParam(param []*smn_pglang.VarDef) string {
	return goitf2lang.ToCppParam(param)
}

func TooCppRet(rets []*smn_pglang.VarDef, pkg, itfName, fName string) string {
	return goitf2lang.TooCppRet(rets, pkg, itfName, fName)
}

func cppClientHead(dir string, itf *smn_pglang.ItfDef) (err error) {
	pkg := itf.Package
	f, err := smn_file.CreateNewFile(dir + ".h")

	if err != nil {
		return err
	}

	defer f.Close()

	writef := func(s string, a ...interface{}) {
		_, err = f.WriteString(fmt.Sprintf(s+"\n", a...))
	}

	writef(`#pragma once
#include<vector>
#include"smncpp/socket_itf.h"
`)

	for _, f := range itf.Functions {
		if len(f.Returns) <= 1 {
			continue
		}

		writef("#include \"pb/rip_%s.pb.h\"", pkg)

		break
	}

	writef(`
namespace clt_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	writef("class %s :public %s::%s{\npublic:\n", itf.Name, pkg, itf.Name)
	defer writef("};\n")

	writef("public:")
	writef("\t%s(const smnet::Conn& c):_c(c) {}", itf.Name)
	writef("\t%s(const smnet::Conn&& c):_c(c) {}", itf.Name)

	for _, f := range itf.Functions {
		writef("\t%s %s(%s)override;\n", TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name, ToCppParam(f.Params))
	}

	return err
}

func cppClientSrc(dir string, itf *smn_pglang.ItfDef) error {
	pkg := itf.Package
	f, err := smn_file.CreateNewFile(dir + ".cpp")

	if err != nil {
		return err
	}

	defer f.Close()

	writef := func(s string, a ...interface{}) {
		_, err = f.WriteString(fmt.Sprintf(s+"\n", a...))
	}

	writef(`
#include<vector>

#include"smncpp/socket_itf.h"
#include "smncpp/socket_mtd.h"
#include "%s.%s.h"
#include "pb/smn_base.pb.h"
#include "pb/smn_dict.pb.h"
`, pkg, itf.Name)

	writef("#include \"pb/rip_%s.pb.h\"", pkg)

	writef(`
namespace clt_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	for _, f := range itf.Functions {
		writef("%s %s::%s(%s){\n", TooCppRet(f.Returns, pkg, itf.Name, f.Name), itf.Name, f.Name, ToCppParam(f.Params))
		writef("\tsmn_base::Call call;")
		writef("\trip_%s::%s_%s_Prm prm;", itf.Package, itf.Name, f.Name)

		for _, f := range f.Params {
			writef("\tprm.set_%s(%s)", strings.ToLower(f.Var), f.Var)
		}

		writef(`	call.set_dict(smn_dict::rip_%s_%s_%s_Prm);`, pkg, itf.Name, f.Name)
		writef(`	call.set_msg(prm.SerializeAsString());
	auto result = smnet::writeString(this->_c, call.SerializeAsString());

	if (result != smnet::ConnStatusSucc){
		throw this->_c.lastError();
	}

	smnet::Bytes retBuff;
	smnet::readLenBytes(this->_c, retBuff);
	smn_base::Ret ret;
`)
		writef(`	rip_%s::%s_%s_Ret lret;	`, pkg, itf.Name, f.Name)
		writef(`	ret.ParseFromArray(retBuff.arr, retBuff.size());
	lret.ParseFromString(ret.msg());`)

		switch len(f.Returns) {
		case 0:
			writef("return;")
		case 1:
			writef("\treturn lret.%s();", strings.ToLower(f.Returns[0].Var))
		default:
			writef("\treturn lret;")
		}

		writef("}")
	}

	return err
}

func CppClient(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	if !smn_file.IsFileExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dir := path + "/" + itf.Package + ".clt." + itf.Name
	err := cppClientHead(dir, itf)

	if err != nil {
		return err
	}

	return cppClientSrc(dir, itf)
}
