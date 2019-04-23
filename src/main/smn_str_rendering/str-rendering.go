package main

import (
	"basis/smn_file"
	"basis/smn_stream"
	"flag"
	"github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
	"os"
	"strings"
	"sync"
	"text/template"
	"unicode"
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

type CallFunc struct {
	name string
	vm   *otto.Otto
}

func getNoNilLine(infile string) map[string]int {
	ck := make(map[string]int)
	fl := strings.Split(string(fileReadAll(infile)), "\n")
	for _, str := range fl {
		str = strings.TrimSpace(str)
		if str == "" {
			continue
		}
		ck[str] = 1
	}
	return ck
}

func (this *CallFunc) Call(params ...interface{}) (otto.Value, error) {
	return this.vm.Call(this.name, nil, params)
}

func getFuncList(funcfile, funclist string) map[string]int {
	fin := smn_stream.FileReadPipeline{FileName: funcfile}
	checkerr(fin.Capture())
	res := make(map[string]int)
	if smn_file.IsFileExist(funclist) {
		res = getNoNilLine(funclist)
	}
	for fin.RemainingSize() != 0 {
		bytes, err := fin.ByteBreakRead('\n')
		if err != nil && err.Error() != "EOF" {
			panic(err)
		}
		str := strings.TrimSpace(string(bytes))
		if strings.HasPrefix(str, "function") {
			str = strings.TrimSpace(str[8:])
			strlen := len(str)
			endIdx := 0
			for ; endIdx < strlen; endIdx++ {
				if !unicode.IsLetter(rune(str[endIdx])) && !unicode.IsNumber(rune(str[endIdx])) && str[endIdx] != '$' && str[endIdx] != '_' {
					break
				}
			}
			str = str[:endIdx]
			res[str] = 1
			checkerr(err)
		}
	}
	return res
}

func loadFuncFile(funcfile string, funclist string) template.FuncMap {
	res := template.FuncMap{}
	bytes := fileReadAll(funcfile)
	vm := otto.New()
	_, err := vm.Run(string(bytes))
	checkerr(err)
	for str := range getFuncList(funcfile, funclist) {
		res[str] = (&CallFunc{name: str, vm: vm}).Call
	}
	return res
}

type citf interface {
	Iadd(k string, v int) int
	Imult(k string, v int) int
	Idiv(k string, v int) int
	Iget(k string) int
	Iset(k string, v int) string
	Inz(k string) bool
}

type Counter struct {
	vs   map[string]int
	lock sync.Mutex
}

func (this *Counter) Iadd(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val += v
	this.vs[k] = val
	return val
}

func (this *Counter) Imult(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val *= v
	this.vs[k] = val
	return v
}

func (this *Counter) Idiv(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val /= v
	this.vs[k] = val
	return v
}

func (this *Counter) Iget(k string) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.vs[k]
}

func (this *Counter) Inz(k string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.vs[k] != 0
}

func (this *Counter) Iset(k string, v int) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.vs[k] = v
	return ""
}

func main() {
	basepath := flag.String("basepath", "./", "base path.")
	input := flag.String("input", "input.tpl", "template file.")
	rdata := flag.String("rdata", "rdata.json", "rendering config data.")
	output := flag.String("output", "output.out", "output file.")
	funcfile := flag.String("funcfile", "func.js", "function js file. not must to configure.")
	funclist := flag.String("funclist", "func.list", "function list. for init")
	flag.Parse()
	*input = *basepath + *input
	*rdata = *basepath + *rdata
	*output = *basepath + *output
	*funcfile = *basepath + *funcfile
	*funclist = *basepath + *funclist
	bytes := fileReadAll(*rdata)
	inp := fileReadAll(*input)
	out := createNewFile(*output)
	t := template.New("zzzz")
	if smn_file.IsFileExist(*funcfile) {
		t.Funcs(loadFuncFile(*funcfile, *funclist))
	}
	ictr := &Counter{vs: make(map[string]int)}
	t.Funcs(template.FuncMap{"iadd": ictr.Iadd, "imult": ictr.Imult, "idiv": ictr.Idiv, "iget": ictr.Iget, "inz": ictr.Inz, "iset": ictr.Iset})
	t, _ = t.Parse(string(inp))
	m := Json2Map(bytes)
	err := t.Execute(&MyWriter{file: out, IsWriteToConsole: true}, m)
	checkerr(err)
}
