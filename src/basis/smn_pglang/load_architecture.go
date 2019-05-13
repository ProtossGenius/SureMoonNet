package smn_pglang

import (
	"basis/smn_data"
	"basis/smn_file"
	"os"
	"strings"
)

func LoadSystem(folderPath, sysName string) (sMap map[string]interface{}, err error) {
	sMap = make(map[string]interface{})
	smn_file.DeepTraversalDir(folderPath, func(path string, info os.FileInfo) bool {
		if !info.IsDir() {
			ts := &System{}
			bytes, err := smn_file.FileReadAll(path)
			if iserr(err) {
				return false
			}
			smn_data.GetDataFromStr(string(bytes), &ts)
			ts.Name = strings.Split(info.Name(), ".")[0]
			sMap[ts.Name] = ts
		}
		return true
	})
	return
}
