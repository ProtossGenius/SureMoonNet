package smn_flag

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
)

type smnFlag struct {
	SFValRegMap map[string]*sFvalReg
}

type sFvalReg struct {
	StrPtr  *string
	BoolPtr *bool
	Func    ActionDo
}

func (this *sFvalReg) GetValue(str string) {
	if this.StrPtr != nil {
		*this.StrPtr = str
	}
	if this.BoolPtr != nil {
		*this.BoolPtr = (str == "true" || str == "t" || str == "")
	}
}

type ActionDo func(args []string) error

func newsmnFlag() *smnFlag {
	res := &smnFlag{SFValRegMap: map[string]*sFvalReg{}}
	return res
}

var _smnFlag = newsmnFlag()

func (this *smnFlag) RegisterString(name string, val *string, useage string, ad ActionDo) {
	flag.StringVar(val, name, *val, useage)
	this.SFValRegMap[name] = &sFvalReg{StrPtr: val, Func: ad}
}

func RegisterString(name string, val *string, useage string, ad ActionDo) {
	_smnFlag.RegisterString(name, val, useage, ad)
}

func (this *smnFlag) RegisterBool(name string, val *bool, useage string, ad ActionDo) {
	flag.BoolVar(val, name, *val, useage)
	this.SFValRegMap[name] = &sFvalReg{BoolPtr: val, Func: ad}
}

func RegisterBool(name string, val *bool, useage string, ad ActionDo) {
	_smnFlag.RegisterBool(name, val, useage, ad)
}

func Parse(args []string, ed *smn_err.ErrDeal) {
	_smnFlag.Parse(flag.Args(), ed)
}

func (this *smnFlag) Parse(args []string, ed *smn_err.ErrDeal) {
	last := ""
	newArgs := make([]string, 0, len(args))
	for _, arg := range args {
		flag := (arg != "-" && arg[0] == '-')
		var val *sFvalReg
		if last != "" {
			val = this.SFValRegMap[last]
			if val != nil {
				if flag {
					val.GetValue("")
				} else {
					val.GetValue(arg)
				}
			}
		}
		if flag {
			if idx := strings.Index(arg, "="); idx != -1 {
				val = this.SFValRegMap[arg[1:idx]]
				if val != nil {
					val.GetValue(arg[idx+1:])
				}
			} else {
				last = arg[1:]
			}
		} else if last != "" {
			newArgs = append(newArgs, arg)
		}

	}
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
		err := valReg.Func(newArgs)
		ed.OnErr(err)
		fmt.Println("deal func ", name, " finish.")
	}
}
