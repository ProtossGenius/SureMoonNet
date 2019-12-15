package godllclz

import (
	"main/smn_code_tools/godllclz/tpl-consts"
)

var oMap = map[string]string{
	"cpp": tpl_consts.TPL_CPP,
}

type GoDllHead struct {
	FuncNew  string   `json:"func_new"`
	FuncDel  string   `json:"func_del"`
	ItfFuncs []string `json:"itf_funcs"`
}

func main() {
	//h := flag.String("h", "./datas/code_tools/godllclz/str-tag-sys-dllmain.h", "input head-file.")
	//l := flag.String("l", "cpp", "output language, such as cpp, java, etc. to add")
	//o := flag.String("o", "./datas/code_tools/godllclz/", "output path. output is o + c + ext.")
	//c := flag.String("c", "StrTagSys", "class name.")
	//i := flag.String("i", "StrTagSysItf", "interface name.")
	//flag.Parse()

}
