package smn_flag

import (
	"flag"
	"fmt"
	"strings"

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

func (this *SFvalReg) GetValue(str string) {
	if this.StrPtr != nil {
		*this.StrPtr = str
	}
	if this.BoolPtr != nil {
		*this.BoolPtr = (str == "true" || str == "t" || str == "")
	}
}

type ActionDo func(sf *SmnFlag, args []string) error

func NewSmnFlag() *SmnFlag {
	res := &SmnFlag{SFValRegMap: map[string]*SFvalReg{}}
	return res
}

func (this *SmnFlag) RegisterString(name string, val *string, useage string, ad ActionDo) {
	flag.StringVar(val, name, *val, useage)
	this.SFValRegMap[name] = &SFvalReg{StrPtr: val, Func: ad}
}

func (this *SmnFlag) RegisterBool(name string, val *bool, useage string, ad ActionDo) {
	flag.BoolVar(val, name, *val, useage)
	this.SFValRegMap[name] = &SFvalReg{BoolPtr: val, Func: ad}
}

func (this *SmnFlag) Parse(args []string, ed *smn_err.ErrDeal) {
	last := ""
	newArgs := make([]string, 0, len(args))
	for _, arg := range args {
		flag := (arg != "-" && arg[0] == '-')
		var val *SFvalReg
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
		err := valReg.Func(this, newArgs)
		ed.OnErr(err)
	}
}
