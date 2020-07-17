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
#include "smncpp/lockm.h"
#include "smncpp/socket_itf.h"

`)

	for inc := range goitf2lang.CppNeedInc(itf, false, false) {
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
	writef("\tstd::mutex                _lock;")
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
#include "smncpp/socket_itf.h"
#include "smncpp/socket_mtd.h"


#include<vector>
`, pkg, itf.Name)

	for inc := range goitf2lang.CppNeedInc(itf, true, true,
		`#include "pb/smn_base.pb.h"`, fmt.Sprintf("#include \"pb/rip_%s.pb.h\"", pkg), `#include "pb/smn_dict.pb.h"`) {
		writef(inc)
	}

	writef(`
namespace clt_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	for _, f := range itf.Functions {
		writef("%s %s::%s(%s){\n", goitf2lang.TooCppRet(f.Returns, pkg, itf.Name, f.Name), itf.Name, f.Name,
			goitf2lang.ToCppParam(f.Params, true))
		writef("\tsmnet::SMLockMgr __s_m_l_o_c_k__(this->_lock);")
		writef("\tsmn_base::Call __s_m_c_a_l_l__;")
		writef("\trip_%s::%s_%s_Prm __s_m_p_r_m__;", itf.Package, itf.Name, f.Name)

		netFunc := ""

		for _, prm := range f.Params {
			if strings.Contains(prm.Type, "net.Conn") {
				netFunc = prm.Var
				continue
			}

			CppFillPb("__s_m_p_r_m__", prm, prm.Var, writef)
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
			writef("\treturn;")
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

func cppServerHead(dir string, itf *smn_pglang.ItfDef) error {
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
#include <memory>
#include <string>
#include <functional>
#include "smncpp/asio_server.h"
#include "smn_itf/%s.%s.h"
`, pkg, itf.Name)

	writef(fmt.Sprintf("#include \"pb/rip_%s.pb.h\"", pkg))
	writef(`
namespace svr_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	writef("class %s:public smnet::Session{\n", itf.Name)
	defer writef("};\n")
	writef("\ttypedef  boost::asio::ip::tcp tcp;")

	writef("public:")
	writef(`	%s(tcp::socket socket, std::shared_ptr<%s::%s> itf):Session(std::move(socket)), _itf(itf),
		 readLen(0), pReadLen(static_cast<char*>((void*)&readLen)){}`, itf.Name, pkg, itf.Name)

	for _, f := range itf.Functions {
		writef("\t%s %s(%s);\n", goitf2lang.TooCppRet(f.Returns, pkg, itf.Name, f.Name), f.Name,
			goitf2lang.ToCppParam(f.Params, true))
	}

	writef("	private:")
	writef(`	void run() override;`)

	for _, f := range itf.Functions {
		writef("\tstd::string %s(const rip_%s::%s_%s_Prm& prm);\n", f.Name,
			pkg, itf.Name, f.Name)
	}

	writef(`	public:		
void pack(const std::string& pb);
	private:
		std::shared_ptr<%s::%s> _itf;
		int32_t readLen;
		char *pReadLen;
`, pkg, itf.Name)

	return err
}

func cppServerSrc(dir string, itf *smn_pglang.ItfDef) error {
	pkg := itf.Package
	f, err := smn_file.CreateNewFile(dir + ".cpp")

	if err != nil {
		return err
	}

	defer f.Close()

	writef := func(s string, a ...interface{}) {
		_, err = f.WriteString(fmt.Sprintf(s+"\n", a...))
	}

	writef(`#include "%s.svr.%s.h"`, pkg, itf.Name)
	writef(`
#include <vector>
#include <iostream>
#include "smncpp/socket_mtd.h"
`)

	for inc := range goitf2lang.CppNeedInc(itf, true, true, `#include "pb/smn_base.pb.h"`, `#include "pb/smn_dict.pb.h"`) {
		writef(inc)
	}

	writef(`
namespace svr_rpc_%s{

`, pkg)

	defer writef("}//namespace clt_rpc_%s", pkg)

	writef(`void %s::pack(const std::string& pb){
			smn_base::Ret ret;
			ret.set_msg(pb);
			smnet::writeString(this->_conn,ret.SerializeAsString());
		}
`, itf.Name)

	writef("void %s::run(){", itf.Name)
	writef(`	auto self = shared_from_this();
	auto& sock = this->_conn->getSocket();
	sock.async_read_some(boost::asio::buffer(pReadLen, 4), [this, self](boost::system::error_code ec, 
			std::size_t lLen){
		if (ec){return;}
		if (lLen < 4){return;}
		smnet::netEdianChange(this->pReadLen, 4);
		smnet::Bytes buff(this->readLen);
		this->_conn->read(this->readLen, buff.arr);
		smn_base::Call call;
		call.ParseFromArray(buff.arr, buff.size());
		try{
		switch(call.dict()){
`)

	for _, f := range itf.Functions {
		writef("\t\tcase smn_dict::rip_%s_%s_%s_Prm:{", pkg, itf.Name, f.Name)
		writef("\t\t\trip_%s::%s_%s_Prm prm;", pkg, itf.Name, f.Name)
		writef("\t\t\tprm.ParseFromString(call.msg());")
		writef("\t\t\tpack(%s(prm));", f.Name)
		writef("\t\t\tbreak;")
		writef("\t\t\t}")
	}

	writef(`		default:{
			std::stringstream ss;
			ss << "error in :" <<  __FILE__ << ":" << __LINE__ <<":func run(), unknow call.dict() = " << call.dict() ;
			throw std::runtime_error(ss.str());
			}
		}
		}catch(std::exception& e){
			std::cout << "Error happened when dealing Call, error is : " << e.what() <<std::endl;
			smn_base::Ret ret;
			ret.set_err(true);
			ret.set_msg(e.what());
			smnet::writeString(this->_conn, ret.SerializeAsString());
		}
		run();
	});
}`) // run;

	for _, f := range itf.Functions {
		writef("std::string %s::%s(const rip_%s::%s_%s_Prm& prm){", itf.Name, f.Name,
			pkg, itf.Name, f.Name)

		prmList := []string{}

		for i, prm := range f.Params {
			cppVarName := strings.ToLower(prm.Var)

			if prm.ArrSize == 0 {
				if strings.Contains(prm.Type, "net.Conn") {
					prmList = append(prmList, "this->_conn")
				} else {
					prmList = append(prmList, fmt.Sprintf("prm.%s()", cppVarName))
				}

				continue
			}

			vd := goitf2lang.ToCppVarDef(prm)
			writef("\t%s __s_m_t_e_m_p_%d__;", vd.Type, i)
			writef("\tfor (int i = 0; i < prm.%s_size(); ++i){__s_m_t_e_m_p_%d__.push_back(prm.%s(i));}", cppVarName,
				i, cppVarName)

			prmList = append(prmList, fmt.Sprintf("__s_m_t_e_m_p_%d__", i))
		}

		if len(f.Returns) == 0 {
			writef("\tthis->_itf->%s(%s);", f.Name, strings.Join(prmList, ", "))
			writef("\treturn \"\";\n}")

			continue
		}

		if len(f.Returns) != 1 {
			writef("\treturn this->_itf->%s(%s).SerializeAsString();\n}", f.Name, strings.Join(prmList, ", "))
			continue
		}

		writef("\tauto itfRet = this->_itf->%s(%s);", f.Name, strings.Join(prmList, ", "))
		writef("\trip_%s::%s_%s_Ret ret;", pkg, itf.Name, f.Name)

		CppFillPb("ret", f.Returns[0], "itfRet", writef)

		writef("\treturn ret.SerializeAsString();")
		writef("}")
	}

	return err
}

//CppFillPb .
func CppFillPb(pbName string, f *smn_pglang.VarDef, varName string, writef func(f string, a ...interface{})) {
	if strings.Contains(f.Type, "net.Conn") {
		return
	}

	fcVar := strings.ToLower(f.Var)

	if f.ArrSize == 0 {
		if strings.Contains(f.Type, ".") {
			vd := goitf2lang.ToCppVarDef(f)
			writef("\t%s.set_allocated_%s(new %s(%s));", pbName, fcVar, vd.Type, varName)
		} else {
			writef("\t%s.set_%s(%s);", pbName, fcVar, varName)
		}

		return
	}

	if !strings.Contains(f.Type, ".") {
		writef("\tfor(size_t i = 0; i < %s.size(); ++i){", varName)

		if f.Type == "string" {
			writef("\t\t%s.add_%s();")
			writef("\t\t%s.set_%s(i, %s[i]);", pbName, fcVar, varName)
		} else {
			writef("\t\t%s.add_%s(%s[i]);", pbName, fcVar, varName)
		}

		writef("\t}")
	} else {
		writef("\tfor(size_t i = 0; i < %s.size(); ++i){%s.add_%s()->CopyFrom(%s[i]);}", f.Var, pbName, fcVar, varName)
	}
}

//CppServer SMNRPC server code product.
func CppServer(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	if !smn_file.IsFileExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	dir := path + "/" + itf.Package + ".svr." + itf.Name
	err := cppServerHead(dir, itf)

	if err != nil {
		return err
	}

	return cppServerSrc(dir, itf)
}
