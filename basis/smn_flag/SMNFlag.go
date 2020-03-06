package smn_flag

import (
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
)

type SmnFlag struct {
	SFValRegMap map[string]*SFvalReg
	ValueMap    map[string]interface{}
}

type SFvalReg struct {
	StrPtr  *string
	BoolPtr *bool
	Func    ActionDo
}

type ActionDo func(sf *SmnFlag, args []string) error

func NewSmnFlag() *SmnFlag {
	res := &SmnFlag{}
	return res
}

func (this *SmnFlag) RegisterString(name string, val *string, ad ActionDo) {
	this.SFValRegMap[name] = &SFvalReg{StrPtr: val, Func: ad}
	this.ValueMap[name] = val
}

func (this *SmnFlag) RegisterBool(name string, val *bool, ad ActionDo) {
	this.SFValRegMap[name] = &SFvalReg{BoolPtr: val, Func: ad}
}

func (this *SmnFlag) Parse(args []string) {
	for name, valReg := range this.SFValRegMap {
		if *(valReg.StrPtr) != "" && *(valReg.BoolPtr) && valReg.Func != nil {
			fmt.Println("dealing funcs .... ", name)
			err := valReg.Func(this, args)
			smn_err.OnErr(err)
		}
	}
}
