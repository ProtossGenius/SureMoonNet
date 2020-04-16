package goitf2lang

import (
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
)

//use checkerr can let code num less.
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

type TypeTrans func(goType string) string

type VarDefTrans func(vd *smn_pglang.VarDef) *smn_pglang.VarDef

type ParamTrans func(param []*smn_pglang.VarDef) string

type ItfTrans func(itf *smn_pglang.ItfDef) string

type FuncWritePkg func(out, pkg string, list []*smn_pglang.ItfDef)

type FuncWriteItf func(out, pkg string, itf *smn_pglang.ItfDef)

var GoItfToLang = map[string]FuncWritePkg{
	"cpp":  WriteCppPkg,
	"java": WriteJavaPkg,
	"go":   writeGoPkg,
}

func writeGoPkg(out, pkg string, list []*smn_pglang.ItfDef) {}

//not exception-safe, it will panic when err happened.
func WriteInterface(lang, out, pkg string, list []*smn_pglang.ItfDef) {
	f, ok := GoItfToLang[lang]
	if !ok {
		fmt.Println("not support such lang")
		return
	}
	f(out, pkg, list)
}
