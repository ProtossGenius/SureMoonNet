package main

import (
	"com.suremoon.net/basis/smn_analysis_go/line_analysis"
	"com.suremoon.net/basis/smn_data"
	"fmt"
	"strings"
	"time"
)

func main() {
	sm := line_analysis.NewGoAnalysis()
	str := `
package hello
type itf interface {
	f(a int,
		b[] int) (int,
		int)
}

type itf2 interface {
	f(a int,
		b int) (int,
		int)
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
		res := <-result
		if res == nil {
			continue
		}
		if res.ProductType() != line_analysis.ProductType_Itf {
			fmt.Println(res.(*line_analysis.GoPkg).Pkg)
			continue
		}
		out := res.(*line_analysis.GoItf)
		str, err := smn_data.ValToJson(out.Result)
		fmt.Printf("zzzzzzzzzzzzzzzzzzzzzzzzzz %s %v\n", str, err)
	}

}
