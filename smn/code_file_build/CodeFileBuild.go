package code_file_build

import (
	"fmt"
	"io"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_muti_write_cache"
)

const (
	//ErrRepeatParse TODO maybe no use.
	ErrRepeatParse = "ErrRepeatParse"
)

//BlockContainer who can create/add block.
type BlockContainer interface {
	AddBlock(format string, a ...interface{}) *CodeBlock
}

//CFBImpFunc write import.
//pkg is package name. for C/C++ is include file name.
//result is all include(import) line
//e.g. In C++: pkg = "cstdio" , out is "#include<cstdio>"
type CFBImpFunc func(pkg string) string

//CppImp import cpp pakcage.
func CppImp(pkg string) string {
	if len(pkg) == 0 {
		return ""
	}

	if pkg[0] == '"' || pkg[0] == '<' {
		return "#include " + pkg
	}

	return "#include <" + pkg + ">"
}

//JavaImp import java pkg.
func JavaImp(pkg string) string {
	return "import " + pkg + ";"
}

//CFBPkgFunc about pkg's declear, for cpp is namespace.
type CFBPkgFunc func(cf *CodeFile, pkg string)

type CodeFile struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	importable map[string]string
	imported   map[string]bool
	importFunc CFBImpFunc
	writer     io.Writer
}

func NewCodeFile(pkg string, w io.Writer, impFunc CFBImpFunc, pkgFunc CFBPkgFunc, comments ...string) *CodeFile {
	res := &CodeFile{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), importable: make(map[string]string), imported: make(map[string]bool), writer: w, importFunc: impFunc}
	pkgFunc(res, pkg)
	res.WriteHeadLine("")
	for _, comment := range comments {
		pre := "//"
		if strings.HasSuffix(comment, "//") {
			pre = ""
		}
		res.WriteHeadLine(pre + comment)
	}
	return res
}

func (this *CodeFile) AddImports(imp map[string]string) {
	for k, v := range imp {
		this.importable[k] = v
	}
}

func (this *CodeFile) _import(pkg string) {
	if this.imported[pkg] {
		return
	}
	this.imported[pkg] = true
	this.WriteHeadLine(this.importFunc(pkg))
}

func (this *CodeFile) Import(str string) bool {
	if val, ok := this.importable[str]; ok {
		this._import(val)
		return true
	}
	this._import(str)
	return false
}

func (this *CodeFile) Imports(imps ...string) {
	for _, val := range imps {
		this.Import(val)
	}
}

func (this *CodeFile) Write(str string) {
	this.Append(smn_muti_write_cache.NewStrCache(str))
}
func (this *CodeFile) WriteLine(str string) {
	this.Write(str + "\n")
}

func (this *CodeFile) AddBlock(format string, a ...interface{}) *CodeBlock {
	f := newCodeBlock(fmt.Sprintf(format, a...), this, 0)
	this.Append(f)
	return f
}

func (this *CodeFile) Output() (int, error) {
	return this.FileMutiWriteCacheItf.Output(this.writer)
}

type CodeBlock struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	father      *CodeFile
	indentation int
}

func newCodeBlock(def string, father *CodeFile, ind int) *CodeBlock {
	res := &CodeBlock{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), father: father, indentation: ind}
	suf := " {"
	if strings.Contains(def, "{") {
		suf = ""
	}
	res.WriteHeadLine(res._addIndentation(def+suf, 0))
	res.WriteTailLine(res._addIndentation("}", 0))
	return res
}

func (this *CodeBlock) Imports(imports ...string) {
	for _, imp := range imports {
		this.father.Import(imp)
	}
}

func (this *CodeBlock) IndentationAdd(n int) {
	this.indentation += n
}

func (this *CodeBlock) _addIndentation(str string, corr int) string {
	space := ""
	for i := 0; i < this.indentation+corr; i++ {
		space += "    "
	}
	str = strings.Replace(str, "\n", "\n"+space, -1)
	return strings.TrimRight(space+str, " ")
}

func (this *CodeBlock) Write(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	this.Append(smn_muti_write_cache.NewStrCache(str))
}

func (this *CodeBlock) WriteToNewLine(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	this.Append(smn_muti_write_cache.NewStrCache(this._addIndentation(str, 1)))
}

func (this *CodeBlock) WriteLine(format string, a ...interface{}) {
	this.WriteToNewLine(format+"\n", a...)
}

func (this *CodeBlock) AddBlock(format string, a ...interface{}) *CodeBlock {
	f := newCodeBlock(fmt.Sprintf(format, a...), this.father, this.indentation+1)
	this.Append(f)
	return f
}
