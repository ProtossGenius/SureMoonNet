package main

import (
	"basis/smn_str_rendering"
	"flag"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	path := flag.String("basepath", "./datas", "basepath")
	flag.Parse()
	render, err := smn_str_rendering.NewStrRender("hello", *path+"/template.tmp")
	checkerr(err)
	checkerr(render.ReadJsFuncs(*path+"/func.js", *path+"/func.list"))
	checkerr(render.ParseData(*path+"/data.json", *path+"/output.txt"))
}
