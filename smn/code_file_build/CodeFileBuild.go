package code_file_build

import (
	"fmt"
	"io"
	"runtime/debug"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_muti_write_cache"
)

const (
	//ErrRepeatParse TODO maybe no use.
	ErrRepeatParse = "ErrRepeatParse"
)

func onerr(err error) {
	if err != nil {
		fmt.Println("Unexcept error ", err, "Stack:\n", string(debug.Stack()))
	}
}

//BlockContainer who can create/add block.
type BlockContainer interface {
	AddBlock(format string, a ...interface{}) *CodeBlock
}

//CFBImpFunc write import.
//pkg is package name. for C/C++ is include file name.
//result is all include(import) line
//e.g. In C++: pkg = "cstdio" , out is "#include<cstdio>"
type CFBImpFunc func(pkg string) string

//JavaImp import java pkg.
func JavaImp(pkg string) string {
	return "import " + pkg + ";"
}

//CFBPkgFunc about pkg's declear, for cpp is namespace.
type CFBPkgFunc func(cf *CodeFile, pkg string)

//CodeFile create a lang's code file.
type CodeFile struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	importable map[string]string
	imported   map[string]bool
	importFunc CFBImpFunc
	writer     io.Writer
}

//NewCodeFile create a new CodeFile.
func NewCodeFile(pkg string, w io.Writer, impFunc CFBImpFunc, pkgFunc CFBPkgFunc, comments ...string) *CodeFile {
	res := &CodeFile{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(),
		importable: make(map[string]string), imported: make(map[string]bool), writer: w, importFunc: impFunc}
	pkgFunc(res, pkg)
	_, err := res.WriteHeadLine("")
	onerr(err)

	for _, comment := range comments {
		pre := "//"
		if strings.HasSuffix(comment, "//") {
			pre = ""
		}

		_, err = res.WriteHeadLine(pre + comment)
		onerr(err)
	}

	return res
}

//AddImports sometimes maybe can't give package's full-path.
func (cf *CodeFile) AddImports(imp map[string]string) {
	for k, v := range imp {
		cf.importable[k] = v
	}
}

func (cf *CodeFile) _import(pkg string) {
	if cf.imported[pkg] {
		return
	}

	cf.imported[pkg] = true
	_, err := cf.WriteHeadLine(cf.importFunc(pkg))
	onerr(err)
}

//Import import package.
func (cf *CodeFile) Import(str string) bool {
	if val, ok := cf.importable[str]; ok {
		cf._import(val)
		return true
	}

	cf._import(str)

	return false
}

//Imports do import.
func (cf *CodeFile) Imports(imps ...string) {
	for _, val := range imps {
		cf.Import(val)
	}
}

//Write add a str cache.
func (cf *CodeFile) Write(str string) {
	cf.Append(smn_muti_write_cache.NewStrCache(str))
}

//WriteLine add on string line.
func (cf *CodeFile) WriteLine(str string) {
	cf.Write(str + "\n")
}

//AddBlock add a block like for{}.
func (cf *CodeFile) AddBlock(format string, a ...interface{}) *CodeBlock {
	f := newCodeBlock(fmt.Sprintf(format, a...), cf, 0)
	cf.Append(f)

	return f
}

//Output write to file.
func (cf *CodeFile) Output() (int, error) {
	return cf.FileMutiWriteCacheItf.Output(cf.writer)
}

//CodeBlock {}.
type CodeBlock struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	father      *CodeFile
	indentation int
	BlockStart  string // start of Block such as for cpp/go/java code is "{"
	BlockEnd    string //end of block, for cpp class is "};"
	BlockDef    string // block def such as in "for(){}" "for()"is def
	ht          bool   //is head/tail writed.
}

func newCodeBlock(def string, father *CodeFile, ind int) *CodeBlock {
	res := &CodeBlock{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(),
		father: father, indentation: ind, BlockDef: def, BlockStart: "", BlockEnd: "}", ht: false}

	if !strings.Contains(def, "{") {
		res.BlockStart = "{"
	}

	return res
}

//Imports import muti package.
func (cb *CodeBlock) Imports(imports ...string) {
	for _, imp := range imports {
		cb.father.Import(imp)
	}
}

//IndentationAdd for code format.
func (cb *CodeBlock) IndentationAdd(n int) {
	cb.indentation += n
}

func (cb *CodeBlock) _addIndentation(str string, corr int) string {
	space := ""
	for i := 0; i < cb.indentation+corr; i++ {
		space += "    "
	}

	str = strings.Replace(str, "\n", "\n"+space, -1)

	return strings.TrimRight(space+str, " ")
}

//Write writef maybe better.
func (cb *CodeBlock) Write(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	cb.Append(smn_muti_write_cache.NewStrCache(str))
}

//WriteToNewLine .
func (cb *CodeBlock) WriteToNewLine(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	cb.Append(smn_muti_write_cache.NewStrCache(cb._addIndentation(str, 1)))
}

//WriteLine .
func (cb *CodeBlock) WriteLine(format string, a ...interface{}) {
	cb.WriteToNewLine(format+"\n", a...)
}

//AddBlock add block.
func (cb *CodeBlock) AddBlock(format string, a ...interface{}) *CodeBlock {
	f := newCodeBlock(fmt.Sprintf(format, a...), cb.father, cb.indentation+1)
	cb.Append(f)

	return f
}

//Output let user can decide block head & tail.
func (cb *CodeBlock) Output(oup io.Writer) (int, error) {
	if !cb.ht {
		_, err := cb.WriteHeadLine(cb.BlockDef + cb.BlockStart)
		onerr(err)
		_, err = cb.WriteTailLine(cb.BlockEnd)
		onerr(err)

		cb.ht = true
	} else {
		fmt.Println("should not call output more than one times.", debug.Stack())
	}

	return cb.FileMutiWriteCacheItf.Output(oup)
}
