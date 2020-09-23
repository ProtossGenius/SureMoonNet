package line_analysis

import (
	"fmt"
	"strings"

	"github.com/ProtossGenius/pglang/snreader"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis_go/smn_anlys_go_tif"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
)

// only for easy analysis.
type LineInput struct {
	snreader.InputItf
	Input string
}

func (this *LineInput) Copy() snreader.InputItf {
	return &LineInput{Input: this.Input}
}

const (
	ProductType_Struct = iota
	ProductType_Itf
	ProductType_Pkg
)

func NewGoAnalysis() *snreader.StateMachine {
	sm := (&snreader.StateMachine{}).Init()
	dftSNR := snreader.NewDftStateNodeReader(sm)
	dftSNR.Register(&GoStructNodeReader{})
	dftSNR.Register(&GoItfNodeReader{})
	dftSNR.Register(&GoPkgNodeReader{})
	return sm
}

//
//func AnalysisFile(path string) ([]snreader.ProductItf, error) {
//
//}

type GoStruct struct {
	snreader.ProductItf
	Result *smn_pglang.StructDef
}

func (*GoStruct) ProductType() int {
	return ProductType_Struct
}

type GoItf struct {
	snreader.ProductItf
	Result *smn_pglang.ItfDef
}

func (*GoItf) ProductType() int {
	return ProductType_Itf
}

type GoPkg struct {
	snreader.ProductItf
	Pkg string
}

func (*GoPkg) ProductType() int {
	return ProductType_Pkg
}

type GoStructNodeReader struct {
	Result  *GoStruct
	reading bool //is start analysis
}

const (
	ErrNotStructHead     = "ErrNotStructHead: %s"
	ErrNotItfHead        = "ErrNotItfHead: %s"
	ErrNotPkgDef         = "ErrNotPkgDef"
	ErrItfEOF            = "ErrItfEOF"
	ErrItfExitWithStash  = "ErrItfExitWithStash: %s"
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

func (this *GoStructNodeReader) Name() string {
	return "GoStructNodeReader"
}

func (this *GoStructNodeReader) PreRead(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
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

func (this *GoStructNodeReader) Read(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = strings.Replace(strings.TrimSpace(in.Input), "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	result := this.Result.Result
	list := smn_str.NotNullSpaceSplit(in.Input)
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
			list := smn_str.NotNullSpaceSplit(in.Input)
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

func (this *GoStructNodeReader) End(stateNode *snreader.StateNode) (isEnd bool, err error) {
	return true, fmt.Errorf(ErrStructUnknowInput, "EOF")
}

func (this *GoStructNodeReader) GetProduct() snreader.ProductItf {
	return this.Result
}

func (this *GoStructNodeReader) Clean() {
	this.reading = false
	this.Result = &GoStruct{Result: &smn_pglang.StructDef{Variables: make([]*smn_pglang.VarDef, 0)}}
}

type GoItfNodeReader struct {
	Result   *GoItf
	stash    string //if last string not end, will stash here.
	reading  bool   //is start analysis
	hasStash bool   //is waiting finish read the line. like : Func(a int, \n  b int)(bool)
}

func (this *GoItfNodeReader) Name() string {
	return "GoItfNodeReader"
}

func (this *GoItfNodeReader) PreRead(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = smn_str.DropLineComment(in.Input)
	in.Input = strings.Replace(in.Input, "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	if !this.reading {
		if !strings.HasPrefix(in.Input, "type") || !strings.Contains(in.Input, "interface") {
			return true, fmt.Errorf(ErrNotItfHead, in.Input)
		}
	}
	return false, nil
}

func (this *GoItfNodeReader) Read(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = smn_str.DropLineComment(in.Input)
	in.Input = strings.Replace(in.Input, "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	result := this.Result.Result
	list := smn_str.NotNullSpaceSplit(in.Input)
	if !this.reading {
		result.Name = list[1]
		this.reading = true
		return false, nil
	} else {
		endFlag := strings.Contains(in.Input, "}")
		in.Input = strings.Replace(in.Input, "}", "", -1)
		if this.hasStash {
			in.Input = this.stash + in.Input
			this.hasStash = false
		}
		if in.Input != "" {
			if strings.HasSuffix(in.Input, ",") {
				this.stash = in.Input
				this.hasStash = true
				return false, nil
			}
			strings.Replace(in.Input, "\n", " ", -1)
			f, err := smn_anlys_go_tif.GetFuncDefFromStr(in.Input)
			if err != nil {
				return true, err
			}
			result.Functions = append(result.Functions, f)

		}
		if endFlag {
			if this.hasStash {
				return true, fmt.Errorf(ErrItfExitWithStash, this.stash)
			}
			return true, nil
		}
		return false, nil
	}
}

func (this *GoItfNodeReader) End(stateNode *snreader.StateNode) (bool, error) {
	return true, fmt.Errorf(ErrItfEOF)
}

func (this *GoItfNodeReader) GetProduct() snreader.ProductItf {
	return this.Result
}

func (this *GoItfNodeReader) Clean() {
	this.reading = false
	this.hasStash = false
	this.Result = &GoItf{Result: smn_pglang.NewItfDefine()}
}

type GoPkgNodeReader struct {
	Result *GoPkg
}

func (this *GoPkgNodeReader) Name() string {
	return "GoPkgNodeReader"
}
func (this *GoPkgNodeReader) PreRead(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
	return false, nil
}

func (this *GoPkgNodeReader) Read(stateNode *snreader.StateNode, input snreader.InputItf) (isEnd bool, err error) {
	in := input.(*LineInput)
	in.Input = smn_str.DropLineComment(in.Input)
	in.Input = strings.Replace(in.Input, "{", "", -1)
	if in.Input == "" {
		return false, nil
	}
	if strings.HasPrefix(in.Input, "package") {
		pkg := strings.Split(in.Input[7:], "/")[0]
		pkg = strings.TrimSpace(pkg)
		this.Result.Pkg = pkg
		return true, nil
	}
	return true, fmt.Errorf(ErrNotPkgDef)
}

func (this *GoPkgNodeReader) End(stateNode *snreader.StateNode) (bool, error) {
	return true, fmt.Errorf(ErrNotPkgDef)
}

func (this *GoPkgNodeReader) GetProduct() snreader.ProductItf {
	return this.Result
}

func (this *GoPkgNodeReader) Clean() {
	this.Result = &GoPkg{}
}
