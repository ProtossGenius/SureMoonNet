package main

import (
	"basis/smn_str"
	"fmt"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	str := "void"
	a, b := smn_str.AnalysisTwoSplitTrim(str, smn_str.CIdentifierJoinEndCheck, smn_str.CIdentifierDropEndCheck)
	fmt.Printf("|%s|%s|", a, b)
}
