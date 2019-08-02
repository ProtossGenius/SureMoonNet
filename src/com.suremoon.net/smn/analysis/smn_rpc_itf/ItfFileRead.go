package smn_rpc_itf

import (
	"com.suremoon.net/basis/smn_analysis_go/line_analysis"
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/basis/smn_str"
	"os"
	"strings"
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

func GetItfListFromDir(path string) (map[string][]*smn_pglang.ItfDef, error) {
	res := make(map[string][]*smn_pglang.ItfDef)
	var err error
	smn_file.DeepTraversalDir(path, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		if strings.HasSuffix(path, ".go") {
			var r []*smn_pglang.ItfDef
			r, err = GetItfList(path)
			if err != nil {
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			if len(r) == 0 {
				return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
			}
			pkg := r[0].Package
			if rList, ok := res[pkg]; ok {
				rList = append(rList, r...)
				res[pkg] = rList
			} else {
				res[pkg] = r
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})
	return res, err
}
