package main

import (
	"basis/smn_file"
	"fmt"
	"github.com/robertkrimen/otto"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	bytes, err := smn_file.FileReadAll("./datas/test.js")
	checkerr(err)
	vm := otto.New()
	_, err = vm.Run(string(bytes))
	checkerr(err)
	data := "hello"
	value, err := vm.Call("test", nil, data)
	checkerr(err)
	fmt.Println(value.String())
}
