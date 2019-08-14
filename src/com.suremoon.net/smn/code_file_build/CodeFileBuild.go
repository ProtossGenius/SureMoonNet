package code_file_build

import (
	"bufio"
	"com.suremoon.net/basis/smn_analysis_go/line_analysis"
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/basis/smn_muti_write_cache"
	"fmt"
	"io"
	"log"
	"os"
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

func (this *GoFile) Imports(imps ...string) {
	for _, val := range imps {
		this.Import(val)
	}
}

func (this *GoFile) Write(str string) {
	this.Append(smn_muti_write_cache.NewStrCache(str))
}
func (this *GoFile) WriteLine(str string) {
	this.Write(str + "\n")
}

func (this *GoFile) AddBlock(format string, a ...interface{}) *GoBlock {
	f := newGoBlock(fmt.Sprintf(format, a...), this, 0)
	this.Append(f)
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

func newGoBlock(def string, fathre *GoFile, ind int) *GoBlock {
	res := &GoBlock{FileMutiWriteCacheItf: smn_muti_write_cache.NewFileMutiWriteCache(), father: fathre, indentation: ind}
	suf := " {"
	if strings.Contains(def, "{") {
		suf = ""
	}
	res.WriteHeadLine(res._addIndentation(def+suf, 0))
	res.WriteTailLine(res._addIndentation("}", 0))
	return res
}

func (this *GoBlock) Imports(imports ...string) {
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

func (this *GoBlock) Write(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	this.Append(smn_muti_write_cache.NewStrCache(str))
}

func (this *GoBlock) WriteToNewLine(format string, a ...interface{}) {
	str := fmt.Sprintf(format, a...)
	this.Append(smn_muti_write_cache.NewStrCache(this._addIndentation(str, 1)))
}

func (this *GoBlock) WriteLine(format string, a ...interface{}) {
	this.WriteToNewLine(format+"\n", a...)
}

func (this *GoBlock) AddBlock(format string, a ...interface{}) *GoBlock {
	f := newGoBlock(fmt.Sprintf(format, a...), this.father, this.indentation+1)
	this.Append(f)
	return f
}

func LocalImportable(path string) map[string]string {
	res := make(map[string]string)
	smn_file.DeepTraversalDir(path, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		scanner, err := smn_file.FileScanner(path)
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
				path = strings.Replace(path, "./src/", "", -1)
				path = path[:strings.LastIndex(path, "/")]
				res[pkg] = path
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return res
}
