package main

import (
	"github.com/ProtossGenius/pglang/snreader"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis_go/line_analysis"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_data"
	"fmt"
	"strings"
	"time"
)

func main() {
	sm := (&snreader.StateMachine{}).Init()
	dftSNR := snreader.NewDftStateNodeReader(sm)
	dftSNR.Register(&line_analysis.GoStructNodeReader{})
	str := `
type Input struct {
	snreader.InputItf
	Input rune
}
type Output struct {
	snreader.ProductItf
	Result int
}
`
	lines := strings.Split(str, "\n")
	for _, line := range lines {
		err := sm.Read(&line_analysis.LineInput{Input: line})
		if err != nil {
			panic(err)
		}
	}
	result := sm.GetResultChan()
	go func() {
		for {
			time.Sleep(1)
		}
	}()
	for {
		for {
			res := <-result
			if res == nil {
				continue
			}
			out := res.(*line_analysis.GoStruct)
			str, err := smn_data.ValToJson(out.Result)
			fmt.Printf("zzzzzzzzzzzzzzzzzzzzzzzzzz %s %v\n", str, err)
		}
	}
}
