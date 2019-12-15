package smn_godll_head

import (
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_pglang"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
	"strings"
)

func analysisTwoSplitTrim(str string) (string, string) {
	return smn_str.AnalysisTwoSplitTrim(str, smn_str.CIdentifierJoinEndCheck, smn_str.CIdentifierDropEndCheck)
}

func getParamsFromStr(prms string) []*smn_pglang.VarDef {
	prmList := strings.Split(prms, ",")
	res := make([]*smn_pglang.VarDef, 0, len(prmList))
	for _, str := range prmList {
		t, v := analysisTwoSplitTrim(str)
		res = append(res, &smn_pglang.VarDef{Var: v, Type: t})
	}

	return res
}

func getFuncDefFromStr(line string) *smn_pglang.FuncDef {
	f := smn_pglang.NewFuncDef()
	ret, another := analysisTwoSplitTrim(line)
	f.Returns = append(f.Returns, &smn_pglang.VarDef{Type: ret})
	name, another := analysisTwoSplitTrim(another)
	f.Name = name
	prmAres := strings.Split(another, ")")
	for i := range prmAres {
		prmAres[i] = strings.TrimSpace(prmAres[i])
	}
	f.Params = getParamsFromStr(prmAres[0][1:])
	return f
}

// only one interface in one file.
func ReadGodllHead(path string) (res []*smn_pglang.FuncDef) {
	fdata, err := smn_file.FileReadAll(path)
	if iserr(err) {
		return
	}
	lines := strings.Split(string(fdata), "\n")
	res = make([]*smn_pglang.FuncDef, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "extern") && strings.HasSuffix(line, ";") {
			_, line = analysisTwoSplitTrim(line)
			res = append(res, getFuncDefFromStr(line))
		}
	}
	return
}
