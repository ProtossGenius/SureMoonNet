package main

import (
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/basis/smn_str"
	"com.suremoon.net/smn/analysis/smn_rpc_itf"
	"com.suremoon.net/smn/code_file_build"
	"flag"
	"fmt"
	"os"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

/** file as:
package xxxx
import(...)

*/

func i64toi(ot, v string) (string, bool) {
	isArr, typ := smn_str.ProtoUseDeal(ot)
	if !strings.Contains(ot, typ) {
		if !isArr {
			if typ[0] == 'i' {
				return fmt.Sprintf("int(%s)", v), false
			} else {
				return fmt.Sprintf("uint(%s)", v), false
			}
		} else {
			if typ[0] == 'i' {
				return fmt.Sprintf("smn_net.Int64ArrToIntArr(%s)", v), true
			} else {
				return fmt.Sprintf("smn_net.UInt64ArrToUIntArr(%s)", v), true
			}
		}
	} else {
		return fmt.Sprintf("%s", v), false
	}
}

func itoi64(ot, v string) (string, bool) {
	isArr, typ := smn_str.ProtoUseDeal(ot)
	if !strings.Contains(ot, typ) {
		if !isArr {
			if typ[0] == 'i' {
				return fmt.Sprintf("int64(%s)", v), false
			} else {
				return fmt.Sprintf("uint64(%s)", v), false
			}
		} else {
			if typ[0] == 'i' {
				return fmt.Sprintf("smn_net.IntArrToInt64Arr(%s)", v), true
			} else {
				return fmt.Sprintf("smn_net.UIntArrToUInt64Arr(%s)", v), true
			}
		}
	} else {
		return fmt.Sprintf("%s", v), false
	}
}

func writeSvrRpcFile(path string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		file, err := smn_file.CreateNewFile(path + itf.Name + ".go")
		check(err)
		gof := code_file_build.NewGoFile("svr_rpc_"+itf.Package, file, "Product by SureMoonNet", "Author: ProtossGenius", "Auto-code should not change.")
		gof.AddImports(code_file_build.LocalImportable("./src"))
		gof.Imports(itf.Package, "github.com/golang/protobuf/proto")
		{ // rpc struct
			b := gof.AddBlock("type SvrRpc%s struct", itf.Name)
			b.WriteLine("itf %s.%s", itf.Package, itf.Name)
			b.WriteLine("dicts []dict.EDict")
			b.Imports("dict")
		}
		{ // new func
			b := gof.AddBlock("func NewSvrRpc%s(itf %s.%s) *SvrRpc%s", itf.Name, itf.Package, itf.Name, itf.Name)
			b.WriteLine("list := make([]dict.EDict, 0)")
			for _, f := range itf.Functions {
				b.WriteLine("list = append(list, dict.EDict_rip_%s_%s_%s_Prm)", itf.Package, itf.Name, f.Name)
			}
			b.WriteLine("return &SvrRpc%s{itf:itf, dicts:list}", itf.Name)
		}
		{ // used message dict
			b := gof.AddBlock("func (this *SvrRpc%s)getEDictList() []dict.EDict", itf.Name)
			b.WriteLine("return this.dicts")
		}
		{ // struct get net-package
			b := gof.AddBlock("func (this *SvrRpc%s)OnMessage(c *base.Call) (_d dict.EDict, _p proto.Message, _e error)", itf.Name)
			b.Imports("base")
			b.Imports("smn_pbr")
			{ // rb = recover func
				b.WriteLine("defer func() {")
				ib := b.AddBlock("if err := recover(); err != nil {")
				ib.IndentationAdd(1)
				ib.WriteLine("_p = nil")
				ib.Imports("fmt")
				ib.WriteLine("_e = fmt.Errorf(\"%%v\", err)")
				b.WriteLine("}()")
			}
			b.WriteLine("m := smn_pbr.GetMsgByDict(c.Msg, c.Dict)")
			sb := b.AddBlock("switch c.Dict") //sb -> switch block
			for _, f := range itf.Functions {
				cb := sb.AddBlock("case dict.EDict_rip_%s_%s_%s_Prm:", itf.Package, itf.Name, f.Name)
				cb.Imports("rip_" + itf.Package)
				cb.WriteLine("_d = dict.EDict_rip_%s_%s_%s_Ret", itf.Package, itf.Name, f.Name)
				cb.WriteLine("msg := m.(*rip_%s.%s_%s_Prm)", itf.Package, itf.Name, f.Name)
				rets := ""
				for i := 0; i < len(f.Returns); i++ {
					if i != 0 {
						rets += ", "
					}
					rets += fmt.Sprintf("p%d", i)
				}
				cb.WriteToNewLine("%s := this.itf.%s(", rets, f.Name)
				for i, r := range f.Params {
					if i != 0 {
						cb.Write(", ")
					}
					pv, usmn := i64toi(r.Type, "msg."+smn_str.InitialsUpper(r.Var))
					if usmn {
						cb.Imports("smn_net")
					}
					cb.Write(pv)
				}
				cb.Write(")\n")
				cb.WriteToNewLine("return _d, &rip_%s.%s_%s_Ret{", itf.Package, itf.Name, f.Name)
				for i, r := range f.Returns {
					if i != 0 {
						cb.Write(", ")
					}
					pv, usmn := itoi64(r.Type, fmt.Sprintf("p%d", i))
					if usmn {
						cb.Imports("smn_net")
					}
					cb.Write("%s:%s", smn_str.InitialsUpper(r.Var), pv)
				}
				cb.WriteLine("}, nil")
			}
			b.WriteLine("return -1, nil, nil")
		}

		gof.Output()
	}
}

func writeClientRpcFile(path string, list []*smn_pglang.ItfDef) {
	for _, itf := range list {
		file, err := smn_file.CreateNewFile(path + itf.Name + ".go")
		check(err)
		gof := code_file_build.NewGoFile("clt_rpc_"+itf.Package, file, "Product by SureMoonNet", "Author: ProtossGenius", "Auto-code should not change.")
		gof.AddImports(code_file_build.LocalImportable("./src"))
		gof.Imports(itf.Package, "github.com/golang/protobuf/proto")
		gof.Imports("rip_" + itf.Package)
		tryImport := func(typ string) {
			_, typ = smn_str.ProtoUseDeal(typ)
			lst := strings.Split(typ, ".")
			if len(lst) != 1 {
				gof.Imports(lst[0])
			}
		}

		{ // rpc struct
			b := gof.AddBlock("type CltRpc%s struct", itf.Name)
			b.WriteLine("%s.%s", itf.Package, itf.Name)
			b.WriteLine("conn smn_net.MessageAdapterItf")
			b.Imports("dict")
		}
		{ // new func
			b := gof.AddBlock("func NewCltRpc%s(conn smn_net.MessageAdapterItf) *CltRpc%s", itf.Name, itf.Name)
			b.Imports("smn_net")
			b.WriteLine("return &CltRpc%s{conn:conn}", itf.Name)
		}
		{ // interface achieve
			for _, f := range itf.Functions {
				prmList := ""
				resList := ""
				rpcPrms := ""
				rpcRes := ""
				for i, prm := range f.Params {
					tryImport(prm.Type)
					if i != 0 {
						prmList += ", "
						rpcPrms += ", "
					}
					prmList += fmt.Sprintf("%s %s", prm.Var, prm.Type)
					pv, usmn := itoi64(prm.Type, prm.Var)
					rpcPrms += fmt.Sprintf("%s:%s", smn_str.InitialsUpper(prm.Var), pv)
					if usmn {
						gof.Imports("smn_net")
					}
				}
				for i, rp := range f.Returns {
					tryImport(rp.Type)
					if i != 0 {
						resList += ", "
						rpcRes += ", "
					}
					resList += rp.Type
					pv, usmn := i64toi(rp.Type, "res."+smn_str.InitialsUpper(rp.Var))
					rpcRes += pv
					if usmn {
						gof.Imports("smn_net")
					}
				}
				b := gof.AddBlock("func (this *CltRpc%s)%s(%s) (%s)", itf.Name, f.Name, prmList, resList)
				b.WriteLine("msg := &rip_%s.%s_%s_Prm{%s}", itf.Package, itf.Name, f.Name, rpcPrms)
				b.WriteLine("this.conn.WriteCall(dict.EDict_rip_%s_%s_%s_Prm, msg)", itf.Package, itf.Name, f.Name)
				b.WriteLine("rm, err := this.conn.ReadRet()")
				b.WriteLine("if err != nil{\n\tpanic(err)\n}")
				b.WriteLine("if rm.Err{\n\tpanic(string(rm.Msg))\n}")
				b.WriteLine("res := &rip_%s.%s_%s_Ret{}", itf.Package, itf.Name, f.Name)
				b.WriteLine("err = proto.Unmarshal(rm.Msg, res)")
				b.WriteLine("if err != nil{\n\tpanic(err)\n}")
				b.WriteLine("return %s", rpcRes)
			}
		}

		gof.Output()
	}
}

func main() {
	i := flag.String("i", "./src/rpc_itf/", "rpc interface dir.")
	o := flag.String("o", "./src/rpc_nitf/", "rpc interface's net accepter, from proto.Message call interface.")
	s := flag.Bool("s", true, "is product server code")
	c := flag.Bool("c", true, "is product client code")
	flag.Parse()
	itfs, err := smn_rpc_itf.GetItfListFromDir(*i)
	check(err)
	for _, list := range itfs {
		if *s {
			op := *o + "/svrrpc/"
			err := os.MkdirAll(op, os.ModePerm)
			check(err)
			writeSvrRpcFile(op, list)
		}
		if *c {
			op := *o + "/clientrpc/"
			err := os.MkdirAll(op, os.ModePerm)
			check(err)
			writeClientRpcFile(op, list)
		}
	}
}
