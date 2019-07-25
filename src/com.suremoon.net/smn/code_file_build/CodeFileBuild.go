package code_file_build

import (
	"com.suremoon.net/basis/smn_muti_write_cache"
	"io"
	"strings"
)

const (
	ErrRepeatParse = "ErrRepeatParse"
)

type GoFile struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	importable map[string]string
	imported   map[string]bool
	writer     io.Writer
}

func NewGoFile(pkg string, w io.Writer, comments ...string) *GoFile {
	res := &GoFile{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), importable: make(map[string]string), imported: make(map[string]bool), writer: w}
	res.WriteHeadLine("package " + pkg)
	res.WriteHeadLine("")
	for _, comment := range comments {
		pre := "//"
		if strings.HasSuffix(comment, "//") {
			pre = ""
		}
		res.WriteHeadLine(pre + comment)
	}
	res.WriteHeadLine("")
	return res
}

func (this *GoFile) AddImports(imp map[string]string) {
	for k, v := range imp {
		this.importable[k] = v
	}
}

func (this *GoFile) _import(pkg string) {
	if this.imported[pkg] {
		return
	}
	this.imported[pkg] = true
	this.WriteHeadLine("import \"" + pkg + "\"")
}

func (this *GoFile) Import(str string) bool {
	if val, ok := this.importable[str]; ok {
		this._import(val)
		return true
	}
	this._import(str)
	return false
}

func (this *GoFile) Write(str string) {
	this.Append(smn_muti_write_cache.NewStrCache(str))
}
func (this *GoFile) WriteLine(str string) {
	this.Write(str + "\n")
}

func (this *GoFile) AddBlock(def string, imports ...string) *GoBlock {
	f := newGoBlock(def, 0, imports...)
	this.Append(f)
	f.father = this
	return f
}

func (this *GoFile) Output() (int, error) {
	return this.FileMutiWriteCacheItf.Output(this.writer)
}

type GoBlock struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	father      *GoFile
	indentation int
}

func newGoBlock(def string, ind int, imports ...string) *GoBlock {
	res := &GoBlock{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), indentation: ind}
	res._imports(imports...)
	suf := " {"
	if strings.Contains(def, "{") {
		suf = ""
	}
	res.WriteHeadLine(res._addIndentation(def+suf, 0))
	res.WriteTailLine(res._addIndentation("}", 0))
	return res
}

func (this *GoBlock) _imports(imports ...string) {
	for _, imp := range imports {
		this.father.Import(imp)
	}
}

func (this *GoBlock) _addIndentation(str string, corr int) string {
	space := ""
	for i := 0; i < this.indentation+corr; i++ {
		space += "    "
	}
	str = strings.Replace(str, "\n", "\n"+space, -1)
	return strings.TrimRight(space+str, " ")
}

func (this *GoBlock) Write(str string, imports ...string) {
	this._imports(imports...)
	this.Append(smn_muti_write_cache.NewStrCache(this._addIndentation(str, 1)))
}

func (this *GoBlock) WriteLine(str string, imports ...string) {
	this.Write(str+"\n", imports...)
}

func (this *GoBlock) AddBlock(def string, imports ...string) *GoBlock {
	f := newGoBlock(def, this.indentation+1, imports...)
	this.Append(f)
	f.father = this.father
	return f
}
