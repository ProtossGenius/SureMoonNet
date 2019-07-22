package main

import (
	"com.suremoon.net/basis/smn_analysis_cpp/smn_godll_head"
	"com.suremoon.net/basis/smn_analysis_go/smn_anlys_go_tif"
	"com.suremoon.net/basis/smn_data"
	"com.suremoon.net/basis/smn_pglang"
	"com.suremoon.net/basis/smn_str"
	"com.suremoon.net/basis/smn_str_rendering"
	"flag"
	"strings"
)

const TPL_DLL_CLS_NAIN = `//product by to-cpp-class
//from smnet.suremoon.com
#ifndef {{.only_def}}
#define {{.only_def}}
#include "{{.head}}"
namespace {{.package}}{
    class {{.name}} {
    int idx;
    public:
        {{.name}}(){
            idx = New{{.name}}();
        }
        virtual ~{{.name}}(){
            Delete{{.name}}(idx);
        }
    public:
        {{range .functions}}{{range .returns}}{{.type}}{{end}} {{.name}}({{iset "first" 1}}{{range .params}}{{if itrue "first"}}{{iset "first" 0}}{{else}}, {{end}}{{.type}} {{.var}}{{end}}){
            {{if .ret}}return {{end}}{{iset "first" 1}}::{{.name}}(idx, {{range .params}} {{if itrue "first"}}{{iset "first" 0}}{{else}}, {{end}}{{.var}}{{end}});
        }
		{{end}}
    };
}
#endif //{{.only_def}}
`

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

func checkerr(err error) {
	if iserr(err) {
		panic(err)
	}
}
func getLastNoZeroStr(strs []string) string {
	val := ""
	for _, str := range strs {
		if strings.TrimSpace(str) == "" {
			continue
		}
		val = str
	}
	return val
}

func main() {
	head := flag.String("head", "./datas/code_tools/godllclz/str-tag-sys-dllmain.h", "dll's head file.")
	goitf := flag.String("goitf", "./datas/code_tools/godllclz/str-tag-sys-itf.go", "goone interface file.")
	console := flag.Bool("console", true, "show in console.")
	o := flag.String("o", "./datas/code_tools/godllclz/StrTagSysItf.h", "output file.")
	flag.Parse()
	dllInfo := smn_godll_head.ReadGodllHead(*head)
	itf, err := smn_anlys_go_tif.ReadGooneItf(*goitf)
	checkerr(err)
	dllInfoMap := make(map[string]*smn_pglang.FuncDef)
	for _, info := range dllInfo {
		dllInfoMap[info.Name] = info
	}
	for _, finfo := range itf.Functions {
		dInfo := dllInfoMap[finfo.Name]
		finfo.Parse()
		finfo.Returns = dInfo.Returns
		for i, vt := range finfo.Params {
			vt.Type = dInfo.Params[i+1].Type
		}
	}
	dMap, err := smn_data.ValToMap(itf)
	checkerr(err)
	rHead := getLastNoZeroStr(strings.Split(*head, "/"))
	rHead = getLastNoZeroStr(strings.Split(rHead, "\\"))
	dMap["head"] = rHead
	dMap["only_def"] = smn_str.GetConstDefine(rHead)
	render, err := smn_str_rendering.NewStrRender("godll", "", TPL_DLL_CLS_NAIN)
	render.IsWriteToConsole = *console
	checkerr(err)
	err = render.ParseData(dMap, *o)
	checkerr(err)
}
