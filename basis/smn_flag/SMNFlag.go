package smn_flag

import (
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
)

type SmnFlag struct {
	SFValRegMap map[string]*SFvalReg
}

type SFvalReg struct {
	StrPtr  *string
	BoolPtr *bool
	Func    ActionDo
}

type ActionDo func(sf *SmnFlag, args []string) error

func NewSmnFlag() *SmnFlag {
	res := &SmnFlag{SFValRegMap: map[string]*SFvalReg{}}
	return res
}

func (this *SmnFlag) RegisterString(name string, val *string, ad ActionDo) {
	this.SFValRegMap[name] = &SFvalReg{StrPtr: val, Func: ad}
}

func (this *SmnFlag) RegisterBool(name string, val *bool, ad ActionDo) {
	this.SFValRegMap[name] = &SFvalReg{BoolPtr: val, Func: ad}
}

func (this *SmnFlag) Parse(args []string, ed *smn_err.ErrDeal) {
	for name, valReg := range this.SFValRegMap {
		if valReg.Func == nil {
			continue
		}
		if valReg.StrPtr != nil {
			if *(valReg.StrPtr) == "" {
				continue
			}
		} else if valReg.BoolPtr != nil {
			if !*(valReg.BoolPtr) {
				continue
			}
		}
		fmt.Println("dealing funcs .... ", name)
		err := valReg.Func(this, args)
		ed.OnErr(err)
	}
}
