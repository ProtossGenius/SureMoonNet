package main

import (
	"com.suremoon.net/basis/smn_analysis"
	"com.suremoon.net/basis/smn_analysis_go/line_analysis"
	"com.suremoon.net/basis/smn_data"
	"fmt"
	"strings"
	"time"
)

func main() {
	sm := (&smn_analysis.StateMachine{}).Init()
	dftSNR := smn_analysis.NewDftStateNodeReader(sm)
	dftSNR.Register(&line_analysis.GoStructNodeReader{})
	str := `
type Input struct {
	smn_analysis.InputItf
	Input rune
}
type Output struct {
	smn_analysis.ProductItf
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
