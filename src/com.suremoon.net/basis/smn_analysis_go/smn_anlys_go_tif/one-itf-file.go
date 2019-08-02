package smn_anlys_go_tif

//for go to dll.

import (
	"com.suremoon.net/basis/smn_file"
	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/basis/smn_str"
	"fmt"
	"strings"
)

func AnalysisTwoSplitTrim(str string) (string, string) {
	return smn_str.AnalysisTwoSplitTrim(str, smn_str.CIdentifierJoinEndCheck, smn_str.CIdentifierDropEndCheck)
}

func GetParamsFromStr(prms string) []*smn_pglang.VarDef {
	prmList := strings.Split(prms, ",")
	res := make([]*smn_pglang.VarDef, 0, len(prmList))
	for _, str := range prmList {
		v, t := AnalysisTwoSplitTrim(str)
		res = append(res, &smn_pglang.VarDef{Var: v, Type: t})
	}
	lastType := ""
	for i := len(res) - 1; i >= 0; i-- {
		if res[i].Type != "" {
			lastType = res[i].Type
		} else {
			res[i].Type = lastType
		}
	}
	if lastType == "" {
		for i := 0; i < len(res); i++ {
			res[i].Type = res[i].Var
			res[i].Var = fmt.Sprintf("p%d", i)
		}
	}
	return res
}

func GetFuncDefFromStr(line string) (*smn_pglang.FuncDef, error) {
	f := &smn_pglang.FuncDef{}
	name, another := AnalysisTwoSplitTrim(line)
	f.Name = name
	if !strings.Contains(another, ")") || !strings.HasPrefix(another, "(") {
		return nil, fmt.Errorf(ErrUnexpectedGoFunctionDefinition, line)
	}
	prmAres := strings.Split(another, ")")
	for i := range prmAres {
		prmAres[i] = strings.TrimSpace(prmAres[i])
	}
	f.Params = GetParamsFromStr(prmAres[0][1:])
	if len(prmAres) > 1 && strings.HasPrefix(prmAres[1], "(") {
		f.Returns = GetParamsFromStr(prmAres[1][1:])
	}
	return f, nil
}

// only one interface in one file.
func ReadGooneItf(path string) (res *smn_pglang.ItfDef, err error) {
	fdata, err := smn_file.FileReadAll(path)
	if iserr(err) {
		return
	}
	lines := strings.Split(string(fdata), "\n")
	res = smn_pglang.NewItfDefine()
	str := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "package") {
			_, line = AnalysisTwoSplitTrim(line) //drop `package`
			line, _ = AnalysisTwoSplitTrim(line) //get packageName
			res.Package = line
			continue
		}
		if strings.HasPrefix(line, "type") {
			_, line = AnalysisTwoSplitTrim(line) //drop `type`
			line, _ = AnalysisTwoSplitTrim(line) //get interfaceName
			res.Name = line
			continue
		}
		if line[0] == '}' {
			break
		}
		if line[0] < 'A' || line[0] > 'Z' {
			continue
		}
		if strings.HasSuffix(line, ",") {
			str += line
			continue
		}
		if len(str) != 0 {
			line = str + line
			str = ""
		}
		fdef, err := GetFuncDefFromStr(line)
		if iserr(err) {
			return nil, err
		}
		res.Functions = append(res.Functions, fdef)
	}
	return
}
