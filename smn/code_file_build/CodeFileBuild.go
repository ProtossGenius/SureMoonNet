package code_file_build

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis_go/line_analysis"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_muti_write_cache"
)

const (
	ErrRepeatParse = "ErrRepeatParse"
)

//pkg is package name. for C/C++ is include file name.
//result is all include(import) line
//e.g. In C++: pkg = "cstdio" , out is "#include<cstdio>"
type CFBImpFunc func(pkg string) string

func CppImp(pkg string) string {
	return "#include <" + pkg + ">"
}

func GoImp(pkg string) string {
	return "import \"" + pkg + "\""
}

func JavaImp(pkg string) string {
	return "import " + pkg + ";"
}

type CodeFile struct {
	smn_muti_write_cache.FileMutiWriteCacheItf
	importable map[string]string
	imported   map[string]bool
	importFunc CFBImpFunc
	writer     io.Writer
}

func NewCodeFile(pkg string, w io.Writer, f CFBImpFunc, comments ...string) *CodeFile {
	res := &CodeFile{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), importable: make(map[string]string), imported: make(map[string]bool), writer: w, importFunc: f}
	res.WriteHeadLine("package " + pkg)
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

func newCodeBlock(def string, fathre *CodeFile, ind int) *CodeBlock {
	res := &CodeBlock{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), father: fathre, indentation: ind}
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

//In order to be compatible with the previous code.
func NewGoFile(pkg string, w io.Writer, comments ...string) *CodeFile {
	return NewCodeFile(pkg, w, GoImp, comments...)
}

func LocalImptTarget(goPath, targetPath string) map[string]string {
	res := make(map[string]string)
	smn_file.DeepTraversalDir(goPath, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		if !strings.HasPrefix(path, targetPath) {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		scanner, file, err := smn_file.FileScanner(path)
		if err != nil {
			log.Printf("LocalImportable DeepTraversalDir path %s, error %s\n", path, err)
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "package") {
				line = strings.Split(line[8:], "/")[0]
				pkg := line_analysis.NotNullSpaceSplit(line)[0]
				path = strings.Replace(path, "\\", "/", -1)
				path = strings.Replace(path, "//", "/", -1)
				path = strings.Replace(path, goPath, "", -1)
				path = path[:strings.LastIndex(path, "/")]
				for strings.HasPrefix(path, "/") {
					path = path[1:]
				}
				pSplit := strings.Split(path, "/")
				psLen := len(pSplit)
				somePath := pkg
				for i := psLen - 2; i >= 0; i-- {
					somePath = pSplit[i] + "/" + somePath
					res[somePath] = path
				}
				res[pkg] = path
				break
			}
		}
		file.Close()
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return res

}

func LocalImportable(goPath string) map[string]string {
	return LocalImptTarget(goPath, "")
}
