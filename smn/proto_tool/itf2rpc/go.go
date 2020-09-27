package itf2rpc

import (
	"errors"
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
	// NetDotConn net.conn .
	NetDotConn = "net.Conn"
	// SmnBase smn_base .
	SmnBase = "github.com/ProtossGenius/SureMoonNet/pb/smn_base"
	// SmnRPC smn_rpc .
	SmnRPC = "github.com/ProtossGenius/SureMoonNet/smn/net_libs/smn_rpc"
	// SmnConnFunc .
	SmnConnFunc = "smn_rpc.ConnFunc"
)

// ErrAsynClientHaveConn .
var ErrAsynClientHaveConn = errors.New("AsynClint not support user-def conn")

/** file as:
package xxxx
import(...)

*/

func anaVarDefs4Go(vds []*smn_pglang.VarDef, prex string, tryImport func(string),
	gof *code_file_build.CodeFile) (prms, prmDefs, rpcInit, rpcVars, connFunc string) {
	join := func(lst []string) string {
		return strings.Join(lst, ", ")
	}
	size := len(vds)
	prmList := make([]string, 0, size)
	prmDefList := make([]string, 0, size)
	rpcInitList := make([]string, 0, size)
	rpcVarList := make([]string, 0, size)

	for i, vd := range vds {
		tryImport(vd.Type)

		varName := vd.Var
		if varName == fmt.Sprintf("p%d", i) {
			varName = fmt.Sprintf("%s%d", prex, i)
		}

		if strings.TrimSpace(vd.Type) != NetDotConn {
			prmDefList = append(prmDefList, fmt.Sprintf("%s %s", varName, vd.Type))
			prmList = append(prmList, vd.Var)
			pv, usmn := goitoi64(vd.Type, varName)
			rpcInitList = append(rpcInitList, fmt.Sprintf("%s:%s", smn_str.InitialsUpper(vd.Var), pv))

			if usmn {
				gof.Imports(SmnRPC)
			}

			pv, usmn = goi64toi(vd.Type, "_res."+smn_str.InitialsUpper(vd.Var))
			rpcVarList = append(rpcVarList, pv)

			if usmn {
				gof.Imports(SmnRPC)
			}
		} else {
			if connFunc != "" {
				fmt.Println("[warning] have muti conn.")

				continue
			}
			prmDefList = append(prmDefList, fmt.Sprintf("%s %s", vd.Var, SmnConnFunc))
			gof.Import(SmnRPC)
			connFunc = vd.Var
		}
	}

	return join(prmList), join(prmDefList), join(rpcInitList), join(rpcVarList), connFunc
}

func anaFuncDef4Go(f *smn_pglang.FuncDef, tryImport func(string), gof *code_file_build.CodeFile) (prmDefList,
	retDefList, rpcPrms, rpcRes, connFunc string, haveConn bool) {
	_, prmDefList, rpcPrms, _, connFunc = anaVarDefs4Go(f.Params, "p", tryImport, gof)
	haveConn = (connFunc != "")
	_, retDefList, _, rpcRes, _ = anaVarDefs4Go(f.Returns, "r", tryImport, gof)

	return
}

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

func gofmt(filePath string) {
	_ = smn_exec.EasyDirExec("./", "gofmt", "-w", filePath)
}

// GoSvr write to go server RPC code.
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

	defer gofmt(filePath)
	defer file.Close()

	gof := code_file_build.NewGoFile("svr_rpc_"+itf.Package, file,
		"Product by SureMoonNet", "Author: ProtossGenius", "Auto-code should not change.")
	gof.Imports(itfFullPkg, "google.golang.org/protobuf/proto")
	{ //  rpc struct
		b := gof.AddBlock("type SvrRpc%s struct", itf.Name)
		b.WriteLine("itf %s.%s", itf.Package, itf.Name)
		b.WriteLine("dicts []smn_dict.EDict")
		b.Imports(module + "/pb/smn_dict")
	}
	{ //  new func
		b := gof.AddBlock("func NewSvrRpc%s(itf %s.%s) *SvrRpc%s", itf.Name, itf.Package, itf.Name, itf.Name)
		b.WriteLine("list := make([]smn_dict.EDict, 0)")
		for _, f := range itf.Functions {
			b.WriteLine("list = append(list, smn_dict.EDict_rip_%s_%s_%s_Prm)", itf.Package, itf.Name, f.Name)
		}
		b.WriteLine("return &SvrRpc%s{itf:itf, dicts:list}", itf.Name)
	}
	{ //  used message dict
		b := gof.AddBlock("func (this *SvrRpc%s)getEDictList() []smn_dict.EDict", itf.Name)
		b.WriteLine("return this.dicts")
	}
	{ // read proto from bytes
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
	{ //  struct get net-package
		b := gof.AddBlock("func (this *SvrRpc%s)OnMessage(c *smn_base.Call, conn net.Conn)"+
			" (_d int32, _p proto.Message, _e error)", itf.Name)
		b.Imports(SmnBase)
		b.Imports("net")
		{ //  rb = recover func
			b.WriteLine("defer func() {")
			ib := b.AddBlock("if err := recover(); err != nil {")
			ib.IndentationAdd(1)
			ib.WriteLine("_p = nil")
			ib.Imports("fmt")
			ib.WriteLine("_e = fmt.Errorf(\"%%v\", err)")
			b.WriteLine("}()")
		}
		sb := b.AddBlock("switch smn_dict.EDict(c.Dict)") // sb -> switch block
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
		cb.WriteLine(`return -1, nil, fmt.Errorf("Can't Found Dict: %%d", c.Dict)`)
	}

	_, err = gof.Output()

	return err
}

func goClient(path, module string, itf *smn_pglang.ItfDef, crtStruct func(*code_file_build.CodeFile), funcDo func(*smn_pglang.FuncDef,
	*code_file_build.CodeFile, func(string)) error) error {
	realPath := path + "/clt_rpc_" + itf.Package
	if !smn_file.IsFileExist(realPath) {
		err := os.MkdirAll(realPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := realPath + "/" + itf.Name + ".go"
	file, err := smn_file.CreateNewFile(filePath)
	//
	if err != nil {
		return err
	}

	defer gofmt(filePath)
	defer file.Close()

	gof := code_file_build.NewGoFile("clt_rpc_"+itf.Package, file, "Product by SureMoonNet",
		"Author: ProtossGenius", "Auto-code should not change.")

	gof.Imports("google.golang.org/protobuf/proto")
	gof.Imports(module + "/pb/rip_" + itf.Package)

	tryImport := func(typ string) {
		_, typ = smn_str.ProtoUseDeal(typ)
		if typ == NetDotConn {
			gof.Imports(SmnRPC)

			return
		}

		lst := strings.Split(typ, ".")

		if len(lst) != 1 {
			gof.Imports(module + "/pb/" + lst[0])
		}
	}

	crtStruct(gof)
	{ //  interface achieve
		for _, f := range itf.Functions {
			if err := funcDo(f, gof, tryImport); err != nil {
				return err
			}
		}
	}

	_, err = gof.Output()

	return err
}

// GoClient interface to go client RPC code.
func GoClient(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	return goClient(path, module, itf, func(gof *code_file_build.CodeFile) {
		{ //  rpc struct
			b := gof.AddBlock("type CltRpc%s struct", itf.Name)
			b.WriteLine("conn smn_rpc.MessageAdapterItf")
			b.WriteLine("lock sync.Mutex")
			b.Imports(module + "/pb/smn_dict")
			b.Imports(SmnRPC)
			b.Imports("sync")
		}
		{ //  new func
			b := gof.AddBlock("func NewCltRpc%s(conn smn_rpc.MessageAdapterItf) *CltRpc%s", itf.Name, itf.Name)
			b.Imports(SmnRPC)
			b.WriteLine("return &CltRpc%s{conn:conn}", itf.Name)
		}
	}, func(f *smn_pglang.FuncDef, gof *code_file_build.CodeFile, tryImport func(string)) error {
		prmList, retDefList, rpcPrms, rpcRes, connFunc, haveConn := anaFuncDef4Go(f, tryImport, gof)

		b := gof.AddBlock("func (this *CltRpc%s)%s(%s) (%s)", itf.Name, f.Name, prmList, retDefList)
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

		return nil
	})
}

// GoAsynClient interface to go client RPC code.
func GoAsynClient(path, module, itfFullPkg string, itf *smn_pglang.ItfDef) error {
	return goClient(path, module, itf, func(gof *code_file_build.CodeFile) {
		gof.Import(SmnBase)
		{ //  rpc struct
			b := gof.AddBlock("type CltRpc%s struct", itf.Name)
			b.WriteLine("conn     smn_rpc.MessageAdapterItf")
			b.WriteLine("onMsg    chan func(*smn_base.Ret)")
			b.WriteLine("sendChan chan proto.Message")
			b.WriteLine("lock     sync.Mutex")
			b.WriteLine("OnErr    smn_err.OnErr")
			b.Imports(module + "/pb/smn_dict")
			b.Imports(SmnRPC)
			b.Imports("sync")
			b.Imports("github.com/ProtossGenius/SureMoonNet/basis/smn_err")
		}
		{ //  new func
			b := gof.AddBlock("func NewCltRpc%s(conn smn_rpc.MessageAdapterItf, cacheSize int) *CltRpc%s", itf.Name, itf.Name)
			b.Imports(SmnRPC)
			b.WriteLine(`res := &CltRpc%s{conn:conn, onMsg:make(chan func(*smn_base.Ret), cacheSize), 
	sendChan: make(chan proto.Message, cacheSize), OnErr: smn_err.DftOnErr}
	go func() {
		for {
				rcvMsg, err := res.conn.ReadRet()
				if err != nil{
					res.OnErr(err)
				}

				f := <-res.onMsg

				go f(rcvMsg)
			}
		}()
		`, itf.Name)
			b.WriteLine(`
	go func() {
		for {
			_bcall := <-res.sendChan
			_bts, _err := proto.Marshal(_bcall)
			res.OnErr(_err)
			_, _err = smn_net.WriteBytes(_bts, res.conn.GetConn())
			res.OnErr(_err)
		}
	}()
`)
			b.WriteLine("return res")
		}
	}, func(f *smn_pglang.FuncDef, gof *code_file_build.CodeFile, tryImport func(string)) error {
		prmList, retDefList, rpcPrms, rpcRes, _, haveConn := anaFuncDef4Go(f, tryImport, gof)

		if haveConn {
			return ErrAsynClientHaveConn
		}

		b := gof.AddBlock("func (this *CltRpc%s)%s(%s, __sm_callback func(%s))", itf.Name, f.Name, prmList, retDefList)
		gof.Import("github.com/ProtossGenius/SureMoonNet/basis/smn_net")
		b.WriteLine(`__onMsg := func(_rm *smn_base.Ret){
	if _rm.Err{
		this.OnErr(errors.New(string(_rm.Msg)))
	}

	_res := &rip_%s.%s_%s_Ret{}
	_err := proto.Unmarshal(_rm.Msg, _res)
	if _err != nil{
		this.OnErr(_err)
	}
	__sm_callback(%s)
}`, itf.Package, itf.Name, f.Name, rpcRes)
		gof.Import("errors")
		gof.Import("google.golang.org/protobuf/proto")
		b.WriteLine("_msg := &rip_%s.%s_%s_Prm{%s}", itf.Package, itf.Name, f.Name, rpcPrms)
		b.WriteLine("_bts, _err := proto.Marshal(_msg)")
		b.WriteLine("this.OnErr(_err)\n")
		b.WriteLine("_bcall := &smn_base.Call{Dict:int32(smn_dict.EDict_rip_%s_%s_%s_Prm), Msg:_bts}",
			itf.Package, itf.Name, f.Name)
		b.WriteLine("this.OnErr(_err)\n")
		b.WriteLine("this.lock.Lock()")
		b.WriteLine("defer this.lock.Unlock()")
		b.WriteLine("this.sendChan <- _bcall")
		b.WriteLine("this.onMsg <- __onMsg")

		return nil
	})
}
