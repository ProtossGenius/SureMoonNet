package main

import (
	"flag"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_str_rendering"
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
	checkerr(render.ParseFileData(*path+"/data.json", *path+"/output.txt"))
}
