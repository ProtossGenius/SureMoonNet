package main

import (
	"fmt"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis"
)

/*
 *   input set(char) {a, b, c}
 *   result set(char) {1, 2, 3}
 *                    rule
---------------------------------------------------------
a b - 1
a c - 2
b - 3

#########################################################
*/

type Input struct {
	Input rune
}

//Copy should give copy to prevent user change it.
func (i *Input) Copy() smn_analysis.InputItf {
	return &Input{Input: i.Input}
}

type Output struct {
	Result int
}

//ProductType result's type. usually should >= 0.
func (o *Output) ProductType() int {
	return o.Result
}

type Type1NodeReader struct {
	Result *Output
	inputs []*Input
}

func (this *Type1NodeReader) Name() string {
	return "Type1NodeReader"
}

func (this *Type1NodeReader) GetProduct() smn_analysis.ProductItf {
	return this.Result
}

func (this *Type1NodeReader) PreRead(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*Input)
	err = fmt.Errorf("type 1 UnExcept Input %c", rInp.Input)
	if len(this.inputs) == 0 && rInp.Input != 'a' {
		return true, err
	} else if len(this.inputs) == 1 && rInp.Input != 'b' {
		return true, err
	}
	return false, nil
}

func (this *Type1NodeReader) Read(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*Input)
	this.inputs = append(this.inputs, rInp)
	if len(this.inputs) == 2 {
		this.Result = &Output{Result: 1}
		return true, nil
	}
	return false, nil
}

func (this *Type1NodeReader) Clean() {
	this.inputs = make([]*Input, 0)
}

type Type2NodeReader struct {
	Result *Output
	inputs []*Input
}

func (this *Type2NodeReader) Name() string {
	return "Type2NodeReader"
}

func (this *Type2NodeReader) GetProduct() smn_analysis.ProductItf {
	return this.Result
}

func (this *Type2NodeReader) PreRead(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*Input)
	err = fmt.Errorf("type 2 UnExcept Input %c", rInp.Input)
	if len(this.inputs) == 0 && rInp.Input != 'a' {
		return true, err

	} else if len(this.inputs) == 1 && rInp.Input != 'c' {
		return true, err
	}
	return false, nil
}

func (this *Type2NodeReader) Read(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*Input)
	this.inputs = append(this.inputs, rInp)
	if len(this.inputs) == 2 {
		this.Result = &Output{Result: 2}
		return true, nil
	}
	return false, nil
}

func (this *Type2NodeReader) Clean() {
	this.inputs = make([]*Input, 0)
}

type Type3NodeReader struct {
	Result *Output
	inputs []*Input
}

func (this *Type3NodeReader) Name() string {
	return "Type3NodeReader"
}

func (this *Type3NodeReader) GetProduct() smn_analysis.ProductItf {
	return this.Result
}

func (this *Type3NodeReader) PreRead(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	rInp := input.(*Input)
	err = fmt.Errorf("type 3 UnExcept Input %c", rInp.Input)
	if rInp.Input != 'b' {
		return true, err
	}
	return false, nil
}

func (this *Type3NodeReader) Read(stateNode *smn_analysis.StateNode, input smn_analysis.InputItf) (isEnd bool, err error) {
	fmt.Println("save Result is .. ", stateNode.Result)
	rInp := input.(*Input)
	this.inputs = append(this.inputs, rInp)
	this.Result = &Output{Result: 3}
	return true, nil
}

func (this *Type3NodeReader) Clean() {
	this.inputs = make([]*Input, 0)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func testDefault() {
	var err error

	sm := (&smn_analysis.StateMachine{}).Init()
	dftSNR := smn_analysis.NewDftStateNodeReader(sm)
	dftSNR.Register(&Type1NodeReader{})
	dftSNR.Register(&Type2NodeReader{})
	dftSNR.Register(&Type3NodeReader{})

	read := func(cs string) {
		for _, c := range cs {
			err = sm.Read(&Input{Input: c})
			check(err)
		}
	}

	read("abacacabbbb")

	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()

	result := sm.GetResultChan()

	for {
		res := <-result
		if res == nil {
			continue
		}

		out := res.(*Output)
		fmt.Println(out.Result)
	}
}

func testList() {
	var err error

	sm := (&smn_analysis.StateMachine{}).Init()
	dftSNR := smn_analysis.NewDftStateNodeReader(sm)
	dftSNR.Register(smn_analysis.NewStateNodeListReader(&Type1NodeReader{}, &Type2NodeReader{}, &Type3NodeReader{}))

	read := func(cs string) {
		for _, c := range cs {
			err = sm.Read(&Input{Input: c})
			check(err)
		}
	}

	read("abacb")

	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()

	result := sm.GetResultChan()

	for {
		res := <-result
		if res == nil {
			continue
		}

		out := res.(*Output)
		fmt.Println(out)
	}
}

func testSelect() {
	var err error

	sm := (&smn_analysis.StateMachine{}).Init()
	dftSNR := smn_analysis.NewDftStateNodeReader(sm)
	dftSNR.Register(smn_analysis.NewStateNodeSelectReader(&Type1NodeReader{}, &Type2NodeReader{}, &Type3NodeReader{}))

	read := func(cs string) {
		for _, c := range cs {
			err = sm.Read(&Input{Input: c})
			check(err)
		}
	}

	read("abacb")

	go func() {
		for {
			time.Sleep(1 * time.Second)
		}
	}()

	result := sm.GetResultChan()

	for {
		res := <-result
		if res == nil {
			continue
		}

		out := res.(*Output)
		fmt.Println(out)
	}
}

func main() {
	testSelect()
}
