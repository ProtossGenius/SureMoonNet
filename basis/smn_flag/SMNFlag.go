package smn_flag

import (
	"flag"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
)

type smnFlag struct {
	SFValRegMap map[string]*sFvalReg
	Args        []string
}

type sFvalReg struct {
	StrPtr  *string
	BoolPtr *bool
	Func    ActionDo
}

func (this *sFvalReg) SetValue(str string) {
	if this.StrPtr != nil {
		*this.StrPtr = str
	}

	if this.BoolPtr != nil {
		*this.BoolPtr = (str == "true" || str == "t" || str == "" || str == "1")
	}
}

func (this *sFvalReg) GetValue() string {
	if this.StrPtr != nil {
		return *this.StrPtr
	}

	if this.BoolPtr != nil && *this.BoolPtr {
		return "true"
	}

	return ""
}

//GetRegValue get registered value.
func GetRegValue(key string) *sFvalReg {
	return _smnFlag.SFValRegMap[key]
}

func FlagArgs() []string {
	return _smnFlag.Args
}

type ActionDo func(val string) error

func newsmnFlag() *smnFlag {
	res := &smnFlag{SFValRegMap: map[string]*sFvalReg{}}
	return res
}

var _smnFlag = newsmnFlag()

func (this *smnFlag) RegisterString(name string, val *string, usage string, ad ActionDo) {
	flag.StringVar(val, name, *val, usage)

	this.SFValRegMap[name] = &sFvalReg{StrPtr: val, Func: ad}
}

func RegisterString(name string, usage string, ad ActionDo) {
	var val string

	_smnFlag.RegisterString(name, &val, usage, ad)
}

func (this *smnFlag) RegisterBool(name string, val *bool, usage string, ad ActionDo) {
	flag.BoolVar(val, name, *val, usage)

	this.SFValRegMap[name] = &sFvalReg{BoolPtr: val, Func: ad}
}

func RegisterBool(name string, usage string, ad ActionDo) {
	var val bool

	_smnFlag.RegisterBool(name, &val, usage, ad)
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
					val.SetValue("")
				} else {
					val.SetValue(arg)
				}
			}
		}

		if flag {
			if idx := strings.Index(arg, "="); idx != -1 {
				val = this.SFValRegMap[arg[1:idx]]
				if val != nil {
					val.SetValue(arg[idx+1:])
				}
			} else {
				last = arg[1:]
			}
		} else if last != "" {
			newArgs = append(newArgs, arg)
		}
	}

	this.Args = newArgs

	for _, valReg := range this.SFValRegMap {
		if valReg.Func == nil || valReg.GetValue() == "" {
			continue
		}

		err := valReg.Func(valReg.GetValue())
		ed.OnErr(err)
	}
}
