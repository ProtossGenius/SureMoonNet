package main

import (
	"github.com/ProtossGenius/SureMoonNet/smn/code_file_build"
	"os"
)

func main() {
	pkgF := code_file_build.NewGoFile("main", os.Stdout, "hello world!!!!!", "what?", "my fun")
	block := pkgF.AddBlock("func main()")

	block2 := block.AddBlock("for i := 0; i < 10; i ++")
	block2.WriteLine("fmt.Println(i)", "fmt")

	block.WriteLine("//wcnm")
	block.WriteLine("fmt.Println(\"wcnm\")")
	block.WriteLine("fmt.Println(\"wcnm\")")
	block.WriteLine("fmt.Println(\"wcnm\")")

	block3 := block2.AddBlock("for i := 0; i < 10; i ++")
	block3.WriteLine("fmt.Println(i)", "fmt")

	pkgF.Output()
}
