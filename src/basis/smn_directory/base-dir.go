package smn_directory

import (
	"io/ioutil"
	"os"
)

type FileDoFunc func(path string, info os.FileInfo) bool //true continue, false end traversal

func DeepTraversalDir(path string, fileDo FileDoFunc) (info os.FileInfo, err error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, info = range dirs {
		fpath := path + PthSep + info.Name()
		if !fileDo(fpath, info) {
			return info, nil
		}
		if info.IsDir() {
			info, err = DeepTraversalDir(fpath, fileDo)
		}
	}
	return
}
