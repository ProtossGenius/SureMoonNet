package itf2rpc

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/goitf2lang"
)

func hasPkg(typ string) (pkg string) {
	if !strings.Contains(typ, ".") {
		return ""
	}

	pkg = strings.Split(typ, ".")[0]
	pkg = strings.ReplaceAll(pkg, "*", "")
	pkg = strings.ReplaceAll(pkg, "[]", "")
	pkg = strings.TrimSpace(pkg)

	return pkg
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
#include <vector>
#include <memory>
#include <string>
#include <functional>
#include "smncpp/socket_itf.h"
`)

	incMap := map[string]bool{}

	for _, f := range itf.Functions {
		if len(f.Returns) > 1 {
			incMap[fmt.Sprintf("#include \"pb/rip_%s.pb.h\"", pkg)] = true
		} else if len(f.Returns) == 1 {
			if pkg := hasPkg(f.Returns[0].Type); pkg != "" {
				incMap[fmt.Sprintf("#include \"pb/%s.pb.h\"", pkg)] = true
			}
		}

		for _, prm := range f.Params {
			if strings.Contains(prm.Type, "net.Conn") {
				continue
			}

			if pkg := hasPkg(prm.Type); pkg != "" {
				incMap[fmt.Sprintf("#include \"pb/%s.pb.h\"", pkg)] = true
			}
		}
	}

	for inc := range incMap {
		writef(inc)
	}

	writef(`
namespace clt_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	writef("class %s {\n", itf.Name)
	defer writef("};\n")

	writef("private:")
	writef("\tstd::shared_ptr<smnet::Conn> _c;")
	writef("public:")
	writef("\t%s(std::shared_ptr<smnet::Conn> c):_c(c) {}", itf.Name)

	for _, f := range itf.Functions {
		writef("\t%s %s(%s);\n", goitf2lang.TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name,
			goitf2lang.ToCppParam(f.Params, true))
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

	writef(`#include "%s.clt.%s.h"
#include"smncpp/socket_itf.h"
#include "smncpp/socket_mtd.h"
#include "pb/smn_base.pb.h"
#include "pb/smn_dict.pb.h"

#include<vector>
`, pkg, itf.Name)

	writef("#include \"pb/rip_%s.pb.h\"", pkg)

	writef(`
namespace clt_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	for _, f := range itf.Functions {
		writef("%s %s::%s(%s){\n", goitf2lang.TooCppRet(f.Returns, pkg, itf.Name, f.Name), itf.Name, f.Name,
			goitf2lang.ToCppParam(f.Params, true))
		writef("\tsmn_base::Call __s_m_c_a_l_l__;")
		writef("\trip_%s::%s_%s_Prm __s_m_p_r_m__;", itf.Package, itf.Name, f.Name)

		netFunc := ""

		for _, f := range f.Params {
			fcVar := strings.ToLower(f.Var)

			if strings.Contains(f.Type, "net.Conn") {
				netFunc = f.Var
				continue
			}

			if f.ArrSize == 0 {
				if strings.Contains(f.Type, ".") {
					vd := goitf2lang.ToCppVarDef(f)
					writef("\t__s_m_p_r_m__.set_allocated_%s(new %s(%s));", fcVar, vd.Type, f.Var)
				} else {
					writef("\t__s_m_p_r_m__.set_%s(%s);", fcVar, f.Var)
				}
			} else {
				if !strings.Contains(f.Type, ".") {
					writef("\tfor(size_t i = 0; i < %s.size(); ++i){__s_m_p_r_m__.set_%s(i, %s[i]);}", f.Var, fcVar, f.Var)
				} else {
					writef("\tfor(size_t i = 0; i < %s.size(); ++i){__s_m_p_r_m__.add_%s()->CopyFrom(%s[i]);}", f.Var, fcVar, f.Var)
				}
			}
		}

		writef(`	__s_m_c_a_l_l__.set_dict(smn_dict::rip_%s_%s_%s_Prm);`, pkg, itf.Name, f.Name)
		writef(`	__s_m_c_a_l_l__.set_msg(__s_m_p_r_m__.SerializeAsString());
	//write param to server.
	if (smnet::writeString(this->_c, __s_m_c_a_l_l__.SerializeAsString()) != smnet::ConnStatusSucc){
		throw this->_c->lastError();
	}
	`)

		if netFunc != "" {
			writef("\tif(%s(this->_c) != smnet::ConnStatusSucc){throw this->_c->lastError();}", netFunc)
		}

		writef(`
	//read server's return
	smnet::Bytes __s_m_r_e_t_B_u_f_f__;
	if (smnet::readLenBytes(this->_c, __s_m_r_e_t_B_u_f_f__)  != smnet::ConnStatusSucc){
		throw this->_c->lastError();
	}

	smn_base::Ret __s_m_r_e_t__;
`)
		writef(`	rip_%s::%s_%s_Ret __s_m_l_r_e_t__;	`, pkg, itf.Name, f.Name)
		writef(`	__s_m_r_e_t__.ParseFromArray(__s_m_r_e_t_B_u_f_f__.arr, __s_m_r_e_t_B_u_f_f__.size());
	__s_m_l_r_e_t__.ParseFromString(__s_m_r_e_t__.msg());`)

		switch len(f.Returns) {
		case 0:
			writef("return;")
		case 1:
			if f.Returns[0].ArrSize == 0 {
				writef("\treturn __s_m_l_r_e_t__.%s();", strings.ToLower(f.Returns[0].Var))
			} else {
				vd := goitf2lang.ToCppVarDef(f.Returns[0])
				writef("\t%s __s_m_r_e_t_a_r_r__(__s_m_l_r_e_t__.%s_size());", vd.Type, vd.Var)
				writef("\tfor (int i = 0; i < __s_m_l_r_e_t__.%s_size(); ++i){__s_m_r_e_t_a_r_r__[i] = __s_m_l_r_e_t__.%s(i);}",
					strings.ToLower(vd.Var), strings.ToLower(vd.Var))
				writef("\treturn __s_m_r_e_t_a_r_r__;")
			}

		default:
			writef("\treturn __s_m_l_r_e_t__;")
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
