package main

import (
	"github.com/json-iterator/go"
	"io/ioutil"
	"os"
	"text/template"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func IsFileExist(fileName string) bool {
	var exist = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func RemoveFileIfExist(fileName string) error {
	if IsFileExist(fileName) {
		return os.Remove(fileName)
	}
	return nil
}

func NewCreateFile(fileName string) *os.File {
	RemoveFileIfExist(fileName)
	res, err := os.Create(fileName)
	checkerr(err)
	return res
}

func Json2Map(bts []byte) map[string]interface{} {
	res := make(map[string]interface{})
	jsoniter.Unmarshal(bts, &res)
	return res
}

func fileReadAll(path string) []byte {
	cfg, err := os.Open(path)
	checkerr(err)
	bytes, err := ioutil.ReadAll(cfg)
	checkerr(err)
	return bytes
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
	out := NewCreateFile("./datas/output.txt")
	t := template.New("zzzz")
	t, _ = t.Parse(string(inp))
	err := t.Execute(&MyWriter{file: out, IsWriteToConsole: true}, Json2Map(bytes))
	checkerr(err)
}
