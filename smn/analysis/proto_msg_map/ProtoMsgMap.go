package proto_msg_map

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str"
)

type ProtoMsgMap struct {
	Pkg    string
	MsgMap map[string]bool //key is message name, bool is can use(if start with `//` can't use.)
}

func NewProtoMsgMap() *ProtoMsgMap {
	return &ProtoMsgMap{MsgMap: make(map[string]bool)}
}

func GetProtoMsgMap(path string) (m *ProtoMsgMap, err error) {
	data, err := smn_file.FileReadAll(path)
	if err != nil {
		return nil, err
	}
	m = NewProtoMsgMap()
	lines := strings.Split(string(data), "\n")
	var pkg string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "package") {
			pkg = strings.Split(line[7:], ";")[0]
			pkg = strings.TrimSpace(pkg)
			m.Pkg = pkg
			continue
		}
		if strings.Contains(line, "message") {
			n := strings.Split(line, "message")[1]
			n = strings.Split(n, "{")[0]
			n = strings.TrimSpace(n)
			if strings.HasPrefix(line, "//") {
				m.MsgMap[n] = false
			} else {
				m.MsgMap[n] = true
			}

		}
	}
	return
}

type DictConst struct {
	Name string
	Id   int
}

type DictConstList []*DictConst

func (this DictConstList) Len() int {
	return len(this)
}

func (this DictConstList) Less(i, j int) bool {
	return this[i].Id < this[j].Id
}

func (this DictConstList) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
func Dict(in string) (list DictConstList, const2Name map[string]string, err error) {
	dictFileName := "smn_dict.proto"
	contMap := make(map[int]bool)
	oldDef := make(map[string]int)
	newDef := make(map[string]int)
	const2Name = make(map[string]string)
	newDef["None"] = 0
	max := 0
	data, err := smn_file.FileReadAll(in + "/" + dictFileName)

	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if !strings.Contains(line, "\"") && strings.Contains(line, "=") {
				line = strings.Replace(line, ";", "", -1)
				spl := strings.Split(line, "=")
				name := strings.Replace(spl[0], "/", "", -1)
				name = strings.TrimSpace(name)
				index := smn_str.ToInt(strings.TrimSpace(spl[1]))
				contMap[index] = true
				oldDef[name] = index
				if index > max {
					max = index
				}
			}
		}
	} else { //if not found dict.proto, the error not nil, but is right.
		err = nil
	}
	if err != nil {
		return nil, nil, err
	}
	nocList := make([]int, 0)
	nocId := 0
	for i := 1; i <= max; i++ {
		if !contMap[i] {
			nocList = append(nocList, i)
		}
	}
	max++
	//get all const define
	smn_file.DeepTraversalDir(in, func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if !info.IsDir() && !strings.Contains(info.Name(), dictFileName) {
			var pm *ProtoMsgMap
			pm, err = GetProtoMsgMap(path)
			if err != nil {
				return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
			}
			for msg, useable := range pm.MsgMap {
				n := pm.Pkg + "_" + msg
				var dictN = n
				if !useable {
					dictN = "//" + n
				}
				if _, ok := newDef[n]; dictN == n && ok {
					err = fmt.Errorf(ErrNameRepeat, n)
					return smn_file.FILE_DO_FUNC_RESULT_STOP_TRAV
				}
				if id, ok := oldDef[n]; ok {
					newDef[dictN] = id
				} else {
					if nocId < len(nocList) {
						newDef[dictN] = nocList[nocId]
						nocId++
					} else {
						newDef[dictN] = max
						max++
					}
				}
				const2Name[n] = pm.Pkg + "." + msg
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

	list = make(DictConstList, 0)
	for key, val := range newDef {
		list = append(list, &DictConst{Name: key, Id: val})
	}
	sort.Sort(list)
	return list, const2Name, nil
}
