package code_file_build

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis_go/line_analysis"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
)

func GoPkg(cf *CodeFile, pkg string) {
	cf.WriteHeadLine("package " + pkg)
}

//GoImp import go pkg.
func GoImp(pkg string) string {
	return "import \"" + pkg + "\""
}

//In order to be compatible with the previous code.
func NewGoFile(pkg string, w io.Writer, comments ...string) *CodeFile {
	return NewCodeFile(pkg, w, GoImp, GoPkg, comments...)
}

func LocalImptTarget(goPath string, targetPaths ...string) map[string]string {
	goPath = smn_str.PathFmt(goPath)
	for i := range targetPaths {
		targetPaths[i] = smn_str.PathFmt(targetPaths[i])
	}
	res := make(map[string]string)
	smn_file.DeepTraversalDir(goPath, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		if len(targetPaths) != 0 {
			isInTarget := false
			for _, targetPath := range targetPaths {
				if strings.HasPrefix(path, targetPath) {
					isInTarget = true
					break
				}
			}
			if !isInTarget {
				return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
			}
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
				path = smn_str.PathFmt(path)
				path = strings.Replace(path, goPath, "", -1)
				path = strings.Replace(path, "\\", "/", -1)
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
