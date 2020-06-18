package main

import (
	"fmt"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/smn/proto_tool/itf2rpc"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("start .. ")

	prm := itf2rpc.ToCppParam([]*smn_pglang.VarDef{{Type: "[]int64", Var: "p0", ArrSize: -1}})
	fmt.Println(prm)
}
