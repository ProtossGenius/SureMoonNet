package main

import "basis/smn_str_rendering"

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	render, err := smn_str_rendering.NewStrRender("hello", "./datas/to_rendering.tmp")
	checkerr(err)
	checkerr(render.ReadJsFuncs("./datas/test.js", "./datas/func.list"))
	checkerr(render.ParseData("./datas/testinp.json", "./datas/output.txt"))
}
