package main

import (
	"basis/smn_file"
	"github.com/json-iterator/go"
	"os"
	"text/template"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func createNewFile(fileName string) *os.File {
	f, e := smn_file.CreateNewFile(fileName)
	checkerr(e)
	return f
}

func fileReadAll(path string) []byte {
	bytes, err := smn_file.FileReadAll(path)
	checkerr(err)
	return bytes
}

func Json2Map(bts []byte) map[string]interface{} {
	res := make(map[string]interface{})
	jsoniter.Unmarshal(bts, &res)
	return res
}

type MyWriter struct {
	IsWriteToConsole bool
	file             *os.File
}

func (this *MyWriter) Write(p []byte) (n int, err error) {
	if this.IsWriteToConsole {
		os.Stdout.Write(p) //d
	}
	return this.file.Write(p)
}

func main() {
	bytes := fileReadAll("./datas/testinp.json")
	inp := fileReadAll("./datas/to_rendering.tmp")
	out := createNewFile("./datas/output.txt")
	t := template.New("zzzz")
	t, _ = t.Parse(string(inp))
	m := Json2Map(bytes)
	m["_"] = "{{"
	err := t.Execute(&MyWriter{file: out, IsWriteToConsole: true}, m)
	checkerr(err)
}
