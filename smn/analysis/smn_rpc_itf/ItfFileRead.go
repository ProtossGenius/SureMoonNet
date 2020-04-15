package smn_rpc_itf

import (
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_analysis_go/line_analysis"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
)

func getPkg(lines []string) string {
	for _, line := range lines {
		line = smn_str.DropLineComment(line)
		if strings.HasPrefix(line, "package") {
			pkg := strings.Split(line[7:], ";")[0]
			pkg = strings.TrimSpace(pkg)
			return pkg
		}
	}
	return ""
}

func GetItfList(path string) ([]*smn_pglang.ItfDef, error) {
	data, err := smn_file.FileReadAll(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	sm := line_analysis.NewGoAnalysis()
	for _, line := range lines {
		l := strings.TrimSpace(line)
		if strings.HasPrefix(l, "import") || strings.HasSuffix(l, "\"") || strings.HasPrefix(l, ")") {
			continue
		}
		err := sm.Read(&line_analysis.LineInput{Input: line})
		if err != nil {
			return nil, err
		}
	}
	sm.End()
	ch := sm.GetResultChan()
	var pkg string
	res := make([]*smn_pglang.ItfDef, 0)
	for {
		pro := <-ch
		if pro.ProductType() == -1 {
			break
		}
		switch pro.ProductType() {
		case line_analysis.ProductType_Pkg:
			pkg = pro.(*line_analysis.GoPkg).Pkg
		case line_analysis.ProductType_Itf:
			res = append(res, pro.(*line_analysis.GoItf).Result)
		}
	}
	for _, itf := range res {
		itf.Package = pkg
	}
	return res, nil
}

//result map : path -> []*smn_pglang.ItfDef
func GetItfListFromDir(path string) (map[string][]*smn_pglang.ItfDef, error) {
	res := make(map[string][]*smn_pglang.ItfDef)
	var err error
	smn_file.DeepTraversalDir(path, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		path = strings.Replace(path, "\\", "/", -1)
		if strings.HasSuffix(path, ".go") {
			var r []*smn_pglang.ItfDef
			r, err = GetItfList(path)
			if err != nil {
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			if len(r) == 0 {
				return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
			}
			pkgPath := path[:strings.LastIndex(path, "/")]
			if rList, ok := res[pkgPath]; ok {
				rList = append(rList, r...)
				res[pkgPath] = rList
			} else {
				res[pkgPath] = r
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return res, err
}
