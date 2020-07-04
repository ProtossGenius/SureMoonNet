package itf2rpc

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_exec"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
	"github.com/ProtossGenius/SureMoonNet/smn/code_file_build"
)

const (
	//NetDotConn net.conn.
	NetDotConn = "net.Conn"
	//SmnBase smn_base.
	SmnBase = "github.com/ProtossGenius/SureMoonNet/pb/smn_base"
	//SmnRPC smn_rpc.
	SmnRPC = "github.com/ProtossGenius/SureMoonNet/smn/net_libs/smn_rpc"
)

/** file as:
package xxxx
import(...)

*/

func goi64toi(ot, v string) (string, bool) {
	isArr, typ := smn_str.ProtoUseDeal(ot)
	if strings.Contains(ot, typ) {
		return v, false
	}

	if !isArr {
		if typ[0] == 'i' {
			return fmt.Sprintf("int(%s)", v), false
		}

		return fmt.Sprintf("uint(%s)", v), false
	}

	if typ[0] == 'i' {
		return fmt.Sprintf("smn_rpc.Int64ArrToIntArr(%s)", v), true
	}

	return fmt.Sprintf("smn_rpc.UInt64ArrToUIntArr(%s)", v), true
}

func goitoi64(ot, v string) (string, bool) {
	isArr, typ := smn_str.ProtoUseDeal(ot)
	if strings.Contains(ot, typ) {
		return v, false
	}

	if !isArr {
		if typ[0] == 'i' {
			return fmt.Sprintf("int64(%s)", v), false
		}

		return fmt.Sprintf("uint64(%s)", v), false
	}

	if typ[0] == 'i' {
		return fmt.Sprintf("smn_rpc.IntArrToInt64Arr(%s)", v), true
	}

	return fmt.Sprintf("smn_rpc.UIntArrToUInt64Arr(%s)", v), true
}

//GoSvr write to go server RPC code.
func GoSvr(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	realPath := path + "/svr_rpc_" + itf.Package
	if !smn_file.IsFileExist(realPath) {
		err := os.MkdirAll(realPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := realPath + "/" + itf.Name + ".go"
	file, err := smn_file.CreateNewFile(filePath)

	if err != nil {
		return err
	}

	defer smn_exec.EasyDirExec("./", "gofmt", "-w", filePath)
	defer file.Close()

	gof := code_file_build.NewGoFile("svr_rpc_"+itf.Package, file,
		"Product by SureMoonNet", "Author: ProtossGenius", "Auto-code should not change.")
	gof.Imports(itfFullPkg, "google.golang.org/protobuf/proto")
	{ // rpc struct
		b := gof.AddBlock("type SvrRpc%s struct", itf.Name)
		b.WriteLine("itf %s.%s", itf.Package, itf.Name)
		b.WriteLine("dicts []smn_dict.EDict")
		b.Imports(module + "/pb/smn_dict")
	}
	{ // new func
		b := gof.AddBlock("func NewSvrRpc%s(itf %s.%s) *SvrRpc%s", itf.Name, itf.Package, itf.Name, itf.Name)
		b.WriteLine("list := make([]smn_dict.EDict, 0)")
		for _, f := range itf.Functions {
			b.WriteLine("list = append(list, smn_dict.EDict_rip_%s_%s_%s_Prm)", itf.Package, itf.Name, f.Name)
		}
		b.WriteLine("return &SvrRpc%s{itf:itf, dicts:list}", itf.Name)
	}
	{ // used message dict
		b := gof.AddBlock("func (this *SvrRpc%s)getEDictList() []smn_dict.EDict", itf.Name)
		b.WriteLine("return this.dicts")
	}
	{ //read proto from bytes
		for _, f := range itf.Functions {
			protoType := fmt.Sprintf("rip_%s.%s_%s_Prm", itf.Package, itf.Name, f.Name)
			b := gof.AddBlock("func ReadEdict_rip_%s_%s_%s_Prm(bytes []byte) *%s",
				itf.Package, itf.Name, f.Name, protoType)
			b.WriteLine("msg := &%s{}", protoType)
			b.WriteLine("err := proto.Unmarshal(bytes, msg)")
			b.WriteLine("if err != nil { panic(err) }")
			b.WriteLine("return msg")
		}
	}
	{ // struct get net-package
		b := gof.AddBlock("func (this *SvrRpc%s)OnMessage(c *smn_base.Call, conn net.Conn)"+
			" (_d int32, _p proto.Message, _e error)", itf.Name)
		b.Imports(SmnBase)
		b.Imports("net")
		{ // rb = recover func
			b.WriteLine("defer func() {")
			ib := b.AddBlock("if err := recover(); err != nil {")
			ib.IndentationAdd(1)
			ib.WriteLine("_p = nil")
			ib.Imports("fmt")
			ib.WriteLine("_e = fmt.Errorf(\"%%v\", err)")
			b.WriteLine("}()")
		}
		sb := b.AddBlock("switch smn_dict.EDict(c.Dict)") //sb -> switch block
		for _, f := range itf.Functions {
			cb := sb.AddBlock("case smn_dict.EDict_rip_%s_%s_%s_Prm:", itf.Package, itf.Name, f.Name)
			cb.WriteLine("_msg := ReadEdict_rip_%s_%s_%s_Prm(c.Msg)", itf.Package, itf.Name, f.Name)
			cb.Imports(module + "/pb/rip_" + itf.Package)
			cb.WriteLine("_d = int32(smn_dict.EDict_rip_%s_%s_%s_Ret)", itf.Package, itf.Name, f.Name)
			rets := ""
			for i := 0; i < len(f.Returns); i++ {
				if i != 0 {
					rets += ", "
				}
				rets += fmt.Sprintf("p%d", i)
			}
			if rets != "" {
				rets += " :="
			}
			cb.WriteToNewLine("%s this.itf.%s(", rets, f.Name)
			for i, r := range f.Params {
				if i != 0 {
					cb.Write(", ")
				}

				if strings.TrimSpace(r.Type) != NetDotConn {
					pv, usmn := goi64toi(r.Type, "_msg."+smn_str.InitialsUpper(r.Var))
					if usmn {
						cb.Imports(SmnRPC)
					}
					cb.Write(pv)
				} else {
					cb.Write("conn")
				}
			}
			cb.Write(")\n")
			cb.WriteToNewLine("return _d, &rip_%s.%s_%s_Ret{", itf.Package, itf.Name, f.Name)
			for i, r := range f.Returns {
				if i != 0 {
					cb.Write(", ")
				}
				pv, usmn := goitoi64(r.Type, fmt.Sprintf("p%d", i))
				if usmn {
					cb.Imports(SmnRPC)
				}
				cb.Write("%s:%s", smn_str.InitialsUpper(r.Var), pv)
			}
			cb.WriteLine("}, nil")
		}
		cb := sb.AddBlock("default:")
		cb.WriteLine(`return -1, nil, fmt.Errorf("Can't Found Dict: %d", c.Dict)`)
	}

	_, err = gof.Output()

	return err
}

//GoClient interface to go client RPC code.
func GoClient(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	realPath := path + "/clt_rpc_" + itf.Package
	if !smn_file.IsFileExist(realPath) {
		err := os.MkdirAll(realPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := realPath + "/" + itf.Name + ".go"
	file, err := smn_file.CreateNewFile(filePath)

	if err != nil {
		return err
	}

	defer smn_exec.EasyDirExec("./", "gofmt", "-w", filePath)
	defer file.Close()

	gof := code_file_build.NewGoFile("clt_rpc_"+itf.Package, file, "Product by SureMoonNet",
		"Author: ProtossGenius", "Auto-code should not change.")

	gof.Imports("google.golang.org/protobuf/proto")
	gof.Imports(module + "/pb/rip_" + itf.Package)

	tryImport := func(typ string) {
		_, typ = smn_str.ProtoUseDeal(typ)
		if typ == NetDotConn {
			return
		}

		lst := strings.Split(typ, ".")

		if len(lst) != 1 {
			gof.Imports(module + "/pb/" + lst[0])
		}
	}

	{ // rpc struct
		b := gof.AddBlock("type CltRpc%s struct", itf.Name)
		b.WriteLine("conn smn_rpc.MessageAdapterItf")
		b.WriteLine("lock sync.Mutex")
		b.Imports(module + "/pb/smn_dict")
		b.Imports(SmnRPC)
		b.Imports("sync")
	}
	{ // new func
		b := gof.AddBlock("func NewCltRpc%s(conn smn_rpc.MessageAdapterItf) *CltRpc%s", itf.Name, itf.Name)
		b.Imports(SmnRPC)
		b.WriteLine("return &CltRpc%s{conn:conn}", itf.Name)
	}
	{ // interface achieve
		for _, f := range itf.Functions {
			prmList := ""
			resList := ""
			rpcPrms := ""
			rpcRes := ""
			connFunc := ""
			haveConn := false
			for i, prm := range f.Params {
				tryImport(prm.Type)
				isConn := strings.TrimSpace(prm.Type) == NetDotConn
				if isConn {
					haveConn = true
				}
				if i != 0 {
					prmList += ", "
					if !isConn {
						rpcPrms += ", "
					}
				}
				if !isConn {
					prmList += fmt.Sprintf("%s %s", prm.Var, prm.Type)
				} else {
					prmList += fmt.Sprintf("%s %s", prm.Var, "smn_rpc.ConnFunc")
					connFunc = prm.Var
					gof.Import(SmnRPC)
				}
				if !isConn {
					pv, usmn := goitoi64(prm.Type, prm.Var)
					rpcPrms += fmt.Sprintf("%s:%s", smn_str.InitialsUpper(prm.Var), pv)
					if usmn {
						gof.Imports(SmnRPC)
					}
				}
			}
			for i, rp := range f.Returns {
				tryImport(rp.Type)
				if i != 0 {
					resList += ", "
					rpcRes += ", "
				}
				resList += rp.Type
				pv, usmn := goi64toi(rp.Type, "_res."+smn_str.InitialsUpper(rp.Var))
				rpcRes += pv
				if usmn {
					gof.Imports(SmnRPC)
				}
			}
			b := gof.AddBlock("func (this *CltRpc%s)%s(%s) (%s)", itf.Name, f.Name, prmList, resList)
			b.WriteLine("this.lock.Lock()")
			b.WriteLine("defer this.lock.Unlock()")
			b.WriteLine("_msg := &rip_%s.%s_%s_Prm{%s}", itf.Package, itf.Name, f.Name, rpcPrms)
			b.WriteLine("this.conn.WriteCall(int32(smn_dict.EDict_rip_%s_%s_%s_Prm), _msg)", itf.Package, itf.Name, f.Name)
			if haveConn {
				b.WriteLine("%s(this.conn.GetConn())", connFunc)
			}
			b.WriteLine("_rm, _err := this.conn.ReadRet()")
			b.WriteLine("if _err != nil{\n\tpanic(_err)\n}")
			b.WriteLine("if _rm.Err{\n\tpanic(string(_rm.Msg))\n}")
			b.WriteLine("_res := &rip_%s.%s_%s_Ret{}", itf.Package, itf.Name, f.Name)
			b.WriteLine("_err = proto.Unmarshal(_rm.Msg, _res)")
			b.WriteLine("if _err != nil{\n\tpanic(_err)\n}")
			b.WriteLine("return %s", rpcRes)
		}
	}

	_, err = gof.Output()

	return err
}

//Go go's rpc.
func Go(itfPath, module string, c, s bool) error {
	return nil
}
