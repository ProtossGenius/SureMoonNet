package smn_file

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
 *  path: file's path
 *	info: file's info
 *  return true continue, false end traversal
 */

type FileDoFuncResult int

const (
	FILE_DO_FUNC_RESULT_DEFAULT FileDoFuncResult = iota
	FILE_DO_FUNC_RESULT_STOP_TRAV
	FILE_DO_FUNC_RESULT_NO_DEAL
)

type FileDoFunc func(path string, info os.FileInfo) FileDoFuncResult

func DeepTraversalDir(path string, fileDo FileDoFunc) (info os.FileInfo, err error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	for _, info = range dirs {
		fpath := path + PthSep + info.Name()
		switch fileDo(fpath, info) {
		case FILE_DO_FUNC_RESULT_STOP_TRAV:
			return info, nil
		case FILE_DO_FUNC_RESULT_NO_DEAL:
			continue
		case FILE_DO_FUNC_RESULT_DEFAULT:
			if info.IsDir() {
				info, err = DeepTraversalDir(fpath, fileDo)
			}
		default:
			continue
		}
	}
	return
}

func Pwd() string {
	dir, _ := filepath.Abs(".")
	return dir
}
