package main

import (
	"basis/smn_data"
	"basis/smn_file"
	"basis/smn_str_rendering"
	"flag"
	"strings"
)

// from DllMain.tmp
const TPL_DLL_NAIN = `package main

// it is product by smnet.suremoon.com

import "C"
import {{.factory_package}} "{{.factory_imp}}"
import {{.interface_pkg}} "{{.interface_imp}}"
import "sync"
var objarr = make([]{{.interface_pkg}}.{{.interface_name}}, 0, 30)
var idx int32 = 0
var oimap = make(map[int32]byte)
var oLock sync.Mutex

//export New{{.interface_name}}
func New{{.interface_name}}() int32{
    oLock.Lock()
    defer oLock.Unlock()
    res := int32(0)
    if len(oimap) != 0{
        for id := range oimap {
            res = id
            delete(oimap, id)
            break
        }
    }else {
        res = idx
        idx++
    }
    obj := pgt_factory.Product{{.interface_name}}()
    if res == int32(len(objarr)){
        objarr = append(objarr, obj)
    }else {
        objarr[res] = obj
    }
    return res
}

//export Delete{{.interface_name}}
func Delete{{.interface_name}}(objid int32) bool {
    if objid >= int32(len(objarr)) || objarr[objid] == nil{
        return false
    }
    objarr[objid] = nil
    delete(oimap, objid)
    return true
}


{{range .func_list}}
//export {{.func_name}}
func {{.func_def}}{
    {{if .ret}}return {{end}}objarr[o_b_j_i_n_d_e_x].{{.func_call}}
}
{{end}}

func main() {
    // Need a main function to make CGO compile package as C shared library
}

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

type FuncInfo struct {
	FuncDef  string `json:"func_def"`
	Ret      bool   `json:"ret"`
	FuncCall string `json:"func_call"`
	FuncName string `json:"func_name"`
}

type DLLMainInfo struct {
	FactoryImp     string     `json:"factory_imp"`
	FactoryPackage string     `json:"factory_package"`
	InterfaceImp   string     `json:"interface_imp"`
	InterfacePkg   string     `json:"interface_pkg"`
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
	info := FuncInfo{Ret: true}
	firstLeftBrackets := strings.Index(def, "(")
	if firstLeftBrackets == -1 || firstLeftBrackets == len(def)-1 {
		panic("when deal def : " + def + " '(' pos error.")
	}
	info.FuncDef = def[:firstLeftBrackets+1] + "o_b_j_i_n_d_e_x int32, " + def[firstLeftBrackets+1:]
	fn__def_ret := strings.Split(def, "(")
	info.FuncName = strings.TrimSpace(fn__def_ret[0])
	info.FuncCall = info.FuncName + " ( "
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
	if len(fn__def_ret) == 2 && strings.TrimSpace(def_ret[1]) == "" {
		info.Ret = false
	}
	return info
}

func main() {
	ipath := flag.String("ipath", "./datas/code_tools/godll/testdata/IPath.go", "interface's path, the file can only have one interface.")
	fimp := flag.String("fimp", "./hello/pgt_factory", "package of interface's factory, can get achieve(use as import fimp in maindll)")
	fpkg := flag.String("fpkg", "", "package of interface's factory, can get achieve")
	iimp := flag.String("iimp", "./hello/pgt_interface", "factory of interface, (use as import iimp in maindll.)")
	ipkg := flag.String("ipkg", "", "package of interface")
	dpath := flag.String("dpath", "./datas/code_tools/godll/testdata/MainDll.go", "dll's main file, will write to there")
	itfn := flag.String("itfn", "Hello", "interface's name, use to check if ipath contains this Interface")
	console := flag.Bool("console", false, "is show reult in console.")
	flag.Parse()

	if *fpkg == "" {
		*fpkg = getLastNoZeroStr(strings.Split(*fimp, "/"))
		*fpkg = getLastNoZeroStr(strings.Split(*fpkg, "\\"))
	}

	if *ipkg == "" {
		*ipkg = getLastNoZeroStr(strings.Split(*iimp, "/"))
		*ipkg = getLastNoZeroStr(strings.Split(*ipkg, "\\"))
	}

	*itfn = strings.TrimSpace(*itfn)
	bytes, err := smn_file.FileReadAll(*ipath)
	checkerr(err)
	ondeal := false
	dmi := &DLLMainInfo{FactoryImp: *fimp, FactoryPackage: *fpkg, InterfaceName: *itfn, InterfacePkg: *ipkg, InterfaceImp: *iimp}
	dmi.FuncList = make([]FuncInfo, 0)
	for _, line := range strings.Split(string(bytes), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "type") && strings.Contains(line, "interface") && strings.Contains(line, *itfn) {
			ondeal = true
			continue
		} else if ondeal && strings.HasPrefix(line, "}") {
			break
		}
		if ondeal {
			if ch := line[0]; ch == '/' || ch == '*' || ch == '_' || (ch >= 'a' && ch < 'z') {
				continue
			}
			info := FuncDefToFuncCall(line)
			dmi.FuncList = append(dmi.FuncList, info)
		}
	}
	dMap, err := smn_data.ValToMap(dmi)
	render, err := smn_str_rendering.NewStrRender("godll", "", TPL_DLL_NAIN)
	render.IsWriteToConsole = *console
	checkerr(err)
	err = render.ParseData(dMap, *dpath)
	checkerr(err)
}
