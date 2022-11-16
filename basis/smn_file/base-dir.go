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
	FILE_DO_FUNC_RESULT_DEFAULT   FileDoFuncResult = iota // continue
	FILE_DO_FUNC_RESULT_STOP_TRAV                         // stop trav
	FILE_DO_FUNC_RESULT_NO_DEAL                           // not deal that file and continue
)

// PathSep path sep .
const PathSep = string(os.PathSeparator)

type FileDoFunc func(path string, info os.FileInfo) FileDoFuncResult

// DeepTraversalDir .
func DeepTraversalDir(path string,
	fileDo func(path string, info os.FileInfo) FileDoFuncResult) (info os.FileInfo, err error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, info = range dirs {
		fpath := path + PathSep + info.Name()
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

// DeepTraversalDirWithSelf .
func DeepTraversalDirWithSelf(path string, fileDo func(path string, info os.FileInfo) FileDoFuncResult) (info os.FileInfo, err error) {
	if self, err := os.Stat(path); err == nil {
		if fileDo(path, self) != FILE_DO_FUNC_RESULT_DEFAULT {
			return self, nil
		}
	} else {
		return self, err
	}

	return DeepTraversalDir(path, fileDo)
}

// ListDirs ..
func ListDirs(path string, dirDo func(dPath string)) (err error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	dirDo(path)

	for _, info := range dirs {
		if !info.IsDir() {
			continue
		}

		err = ListDirs(path+PathSep+info.Name(), dirDo)
		if err != nil {
			return err
		}
	}

	return nil
}

// Pwd like system pwd.
func Pwd() string {
	dir, _ := filepath.Abs(".")

	return dir
}
