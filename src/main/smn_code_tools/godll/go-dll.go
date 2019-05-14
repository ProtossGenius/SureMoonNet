package main

import (
	"basis/smn_data"
	"basis/smn_file"
	"basis/smn_str_rendering"
	"flag"
	"strings"
)

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

type FuncInfo struct {
	FuncDef  string `json:"func_def"`
	Ret      bool   `json:"ret"`
	FuncCall string `json:"func_call"`
	FuncName string `json:"func_name"`
}

type DLLMainInfo struct {
	FactoryPath    string     `json:"factory_path"`
	FactoryPackage string     `json:"factory_package"`
	InterfaceName  string     `json:"interface_name"`
	FuncList       []FuncInfo `json:"func_list"`
}

func getFirstNoZeroStr(strs []string) string {
	for _, str := range strs {
		if strings.TrimSpace(str) == "" {
			continue
		}
		return str
	}
	return ""
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

func FuncDefToFuncCall(def string) FuncInfo {
	info := FuncInfo{FuncDef: def, Ret: true}
	fn__def_ret := strings.Split(def, "(")
	info.FuncName = strings.TrimSpace(fn__def_ret[0])
	info.FuncCall = info.FuncName + " ("
	def_ret := strings.Split(fn__def_ret[1], ")")
	prmds := strings.Split(def_ret[0], ",")
	for idx, val := range prmds {
		c_ddd := strings.Contains(val, "...")
		spl := strings.Split(val, " ")
		val = getFirstNoZeroStr(spl)
		val = strings.TrimSpace(val)
		if idx != 0 {
			info.FuncCall += ", "
		}
		info.FuncCall += val
		if c_ddd && !strings.HasPrefix(val, "...") {
			info.FuncCall += "..."
		}
	}
	info.FuncCall += ")"
	if strings.TrimSpace(def_ret[1]) == "" {
		info.Ret = false
	}
	return info
}

// from DllMain.tmp
const DLL_NAIN_TPL = `package main

// it is product by smnet.suremoon.com

import "C"
import {{.factory_package}} "{{.factory_path}}"

var val{{.interface_name}} = {{.factory_package}}.Product{{.interface_name}}()
{{range .func_list}}
//export {{.func_name}}
func {{.func_def}}{
    {{if .ret}}return {{end}}val{{$.interface_name}}.{{.func_call}}
}
{{end}}
`

func main() {
	ipath := flag.String("ipath", "./datas/code_tools/godll/testdata/IPath.go", "interface's path, the file can only have one interface.")
	fpath := flag.String("fpath", "./hello/pgt_factory", "package of interface's factory, can get achieve")
	fpkg := flag.String("fpkg", "", "package of interface's factory, can get achieve")
	dpath := flag.String("dpath", "./datas/code_tools/godll/testdata/MainDll.go", "dll's main file, will write to there")
	itfn := flag.String("itfn", "Hello", "interface's name, use to check if ipath contains this Interface")
	flag.Parse()

	if *fpkg == "" {
		*fpkg = getLastNoZeroStr(strings.Split(*fpath, "/"))
		*fpkg = getLastNoZeroStr(strings.Split(*fpkg, "\\"))
	}
	*itfn = strings.TrimSpace(*itfn)
	bytes, err := smn_file.FileReadAll(*ipath)
	checkerr(err)
	ondeal := false
	dmi := &DLLMainInfo{FactoryPath: *fpath, FactoryPackage: *fpkg, InterfaceName: *itfn}
	dmi.FuncList = make([]FuncInfo, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "type") && strings.Contains(line, "interface") && strings.Contains(line, *itfn) {
			ondeal = true
			continue
		} else if ondeal && strings.HasPrefix(line, "}") {
			break
		}
		if ondeal {
			info := FuncDefToFuncCall(line)
			dmi.FuncList = append(dmi.FuncList, info)
		}
	}
	dMap, err := smn_data.ValToMap(dmi)
	render, err := smn_str_rendering.NewStrRender("godll", "", DLL_NAIN_TPL)
	render.IsWriteToConsole = true
	checkerr(err)
	err = render.ParseData(dMap, *dpath)
	checkerr(err)
}
