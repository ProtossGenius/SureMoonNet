package lw_analysis

import (
	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"fmt"
	"strings"
)

const (
	PRODUCT_TYPE_STRUCT = iota
	PRODUCT_TYPE_
)

const (
	INPUT_TYPE_IDENTIFIER = iota
	INPUT_TYPE_SYMBOL
	INPUT_TYPE_
)

type LangInput struct {
	smn_analysis.InputItf
	Word string
	Type int
}

type ResultStruct struct {
	Result *smn_pglang.StructDef
}

func (this *ResultStruct) ProductType() int {
	return PRODUCT_TYPE_STRUCT
}

func NewResultStruct() *ResultStruct {
	return &ResultStruct{Result: &smn_pglang.StructDef{Variables: make([]*smn_pglang.VarDef, 0)}}
}

type StructReadNode struct {
	Result         *ResultStruct
	waitStructName bool
	waitVarName    bool //is waiting var name.
}

func (this *StructReadNode) GetProduct() smn_analysis.ProductItf {
	return this.Result
}

func (this *StructReadNode) PreRead(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*LangInput)
	if rInp.Word == "struct" && !this.waitStructName {
		return true, nil
	}
	return false, nil
}

func (this *StructReadNode) Read(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*LangInput)
	result := this.Result.Result
	if rInp.Word == "struct" {
		return false, nil
	}
	if rInp.Word == "#" {
		if this.waitVarName {
			err = fmt.Errorf(ErrWaitVarName)
		}
		return true, err
	}

	if this.waitStructName {
		result.Name = rInp.Word
		this.waitStructName = false
	} else if this.waitVarName {
		result.Variables[len(result.Variables)-1].Var = rInp.Word
		this.waitVarName = false
	} else {
		node := &smn_pglang.VarDef{}
		typ := rInp.Word
		if strings.Contains(typ, "[]") {
			node.ArrSize = -1
			typ = strings.Replace(typ, "[]", "", -1)
		}
		node.Type = typ
		result.Variables = append(result.Variables, node)
		this.waitVarName = true
	}
	return false, nil
}

func (this *StructReadNode) Clean() {
	this.waitStructName = true
	this.waitVarName = false
	this.Result = NewResultStruct()
}

func GetStructStateMachine() *smn_analysis.StateMachine {
	sm := (&smn_analysis.StateMachine{}).Init()
	dftSNR := smn_analysis.NewDftStateNodeReader(sm)
	dftSNR.Register(&StructReadNode{})
	return sm
}
