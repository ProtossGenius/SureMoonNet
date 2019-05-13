package smn_str_rendering

import (
	"basis/smn_data"
	"basis/smn_file"
	"basis/smn_func"
	"fmt"
	"os"
	"strings"
	"text/template"
	"unicode"
)

type CFWriter struct {
	IsWriteToConsole bool
	file             *os.File
}

func (this *CFWriter) Write(p []byte) (n int, err error) {
	if this.IsWriteToConsole {
		os.Stdout.Write(p) //d
	}
	return this.file.Write(p)
}

type StrRender struct {
	tplData          string // template file content
	tpl              *template.Template
	err              error
	IsWriteToConsole bool
	IsParsed         bool
}

func (this *StrRender) noError() bool {
	return this.err == nil
}

func (this *StrRender) inError() bool {
	return this.err != nil
}

func (this *StrRender) readFile(path string) []byte {
	Nil := make([]byte, 0)
	if this.err != nil {
		return Nil
	}
	bytes, err := smn_file.FileReadAll(path)
	if this.noError() {
		return bytes
	}
	this.err = err
	return Nil
}

func (this *StrRender) readTplFile(tplFile string) *StrRender {
	this.tplData = string(this.readFile(tplFile))
	return this
}

func (this *StrRender) ParseFileData(dataFile, outFile string) error {
	bytes := this.readFile(dataFile)
	if this.inError() {
		return this.err
	}
	dataMap, err := smn_data.GetDataMapFromStr(string(bytes))
	if iserr(err) {
		return err
	}
	return this.ParseData(dataMap, outFile)
}
func (this *StrRender) ParseData(data interface{}, outFile string) error {
	if !this.IsParsed {
		if this.inError() {
			return this.err
		}
		this.tpl, this.err = this.tpl.Parse(this.tplData)
		if this.inError() {
			return this.err
		}
	}
	out, err := smn_file.CreateNewFile(outFile)
	if iserr(err) {
		return err
	}
	return this.tpl.Execute(&CFWriter{file: out, IsWriteToConsole: this.IsWriteToConsole}, data)
}

func (this *StrRender) getFuncList(jsContent, funcList string) (res map[string]int) {
	res = make(map[string]int)
	//read from funcList
	if smn_file.IsFileExist(funcList) {
		for _, str := range strings.Split(string(this.readFile(funcList)), "\n") {
			if str = strings.TrimSpace(str); str != "" {
				res[str] = 1
			}
		}
	}
	for _, str := range strings.Split(jsContent, "\n") {
		if str = strings.TrimSpace(str); strings.HasPrefix(str, "function") {
			str = strings.TrimSpace(str[8:])
			strlen := len(str)
			endIdx := 0
			dflag := false // contains $
			for ; endIdx < strlen; endIdx++ {
				if !unicode.IsLetter(rune(str[endIdx])) && !unicode.IsNumber(rune(str[endIdx])) && str[endIdx] != '_' {
					break
				}
				if str[endIdx] == '$' {
					dflag = true
					break
				}
			}
			if !dflag {
				str = str[:endIdx]
				res[str] = 1
			}
		}
	}
	return res
}

func (this *StrRender) ReadJsFuncs(jsPath, funcList string) error {
	if !smn_file.IsFileExist(jsPath) {
		return fmt.Errorf(ERR_FILE_NOT_FOUND, jsPath)
	}
	jsByte := this.readFile(jsPath)
	if this.inError() {
		return this.err
	}
	jff, err := smn_func.NewJsFuncFactory(string(jsByte))
	if iserr(err) {
		return err
	}
	fl := this.getFuncList(string(jsByte), funcList)
	fMap := template.FuncMap{}
	for nm := range fl {
		fMap[nm] = jff.ProductFunc(nm)
	}
	this.tpl.Funcs(fMap)
	return nil
}

func NewStrRender(name, tplFile string) (res *StrRender, err error) {
	res = &StrRender{IsWriteToConsole: true, IsParsed: false}
	tpl := template.New(name)
	ictr := &Counter{vs: make(map[string]int)}
	//built in function.
	tpl.Funcs(template.FuncMap{"iadd": ictr.Iadd, "imult": ictr.Imult, "idiv": ictr.Idiv, "Iset": ictr.Iset, "itrue": ictr.Inz, "iset": ictr.Iset,
		"isadd": ictr.Isadd, "ismult": ictr.Ismult, "isdiv": ictr.Isdiv})
	res.readTplFile(tplFile)
	res.tpl = tpl
	return res, res.err
}
