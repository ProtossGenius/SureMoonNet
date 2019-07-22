package line_analysis

import (
	"com.suremoon.net/basis/smn_analysis"
	"com.suremoon.net/basis/smn_pglang"
	"fmt"
	"strings"
)

// only for easy analysis.
type LineInput struct {
	smn_analysis.InputItf
	Input string
}

type GoStruct struct {
	smn_analysis.ProductItf
	Result *smn_pglang.StructDef
}

type GoStructNodeReader struct {
	Result  *GoStruct
	reading bool
}

const (
	ErrNotStructHead     = "ErrNotStructHead: %s"
	ErrNotStructBody     = "ErrNotStructBody: %s"
	ErrStructUnknowInput = "ErrStructUnknowInput: %s"
)

func NotNullSpaceSplit(inp string) []string {
	strs := strings.Split(inp, " ")
	res := make([]string, 0, len(strs))
	for _, str := range strs {
		if str != "" {
			res = append(res, str)
		}
	}
	return res
}

func (this *GoStructNodeReader) PreRead(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = strings.Replace(strings.TrimSpace(in.Input), "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	if !this.reading {
		if !strings.HasPrefix(in.Input, "type") || !strings.Contains(in.Input, "struct") {
			return true, fmt.Errorf(ErrNotStructHead, in.Input)
		}
	}
	return false, nil
}

func (this *GoStructNodeReader) Read(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = strings.Replace(strings.TrimSpace(in.Input), "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	result := this.Result.Result
	list := NotNullSpaceSplit(in.Input)
	if !this.reading {
		result.Name = list[1]
		this.reading = true
		return false, nil
	} else {
		if strings.HasPrefix(in.Input, "/") || strings.HasPrefix(in.Input, "*") {
			return false, nil
		}
		spl := strings.Split(in.Input, "//")
		in.Input = spl[0]
		endFlag := strings.Contains(in.Input, "}")
		in.Input = strings.Replace(in.Input, "}", "", -1)
		if in.Input != "" {
			list := NotNullSpaceSplit(in.Input)
			var v *smn_pglang.VarDef
			if len(list) < 2 {
				v = &smn_pglang.VarDef{Var: "", Type: list[0]}
			} else {
				v = &smn_pglang.VarDef{Var: list[0], Type: list[1]}
			}
			if strings.Contains(in.Input, "[]") {
				v.ArrSize = -1
			}
			result.Variables = append(result.Variables, v)
		}
		if endFlag {
			return true, nil
		}
		return false, nil
	}
	return true, fmt.Errorf(ErrStructUnknowInput, in.Input)
}

func (this *GoStructNodeReader) GetProduct() smn_analysis.ProductItf {
	return this.Result
}

func (this *GoStructNodeReader) Clean() {
	this.reading = false
	this.Result = &GoStruct{Result: &smn_pglang.StructDef{Variables: make([]*smn_pglang.VarDef, 0)}}
}