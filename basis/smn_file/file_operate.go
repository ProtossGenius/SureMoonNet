package smn_file

import (
	"bufio"
	"io/ioutil"
	"os"
)

//IsFileExist 判断文件是否存在  存在返回 true 不存在返回false.
func IsFileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}

func RemoveFileIfExist(fileName string) error {
	if IsFileExist(fileName) {
		return os.Remove(fileName)
	}
	return nil
}

func SafeOpenFile(fileName string) (*os.File, error) {
	if IsFileExist(fileName) {
		return os.Open(fileName)
	}
	return os.Create(fileName)
}

func CreateNewFile(fileName string) (*os.File, error) {
	RemoveFileIfExist(fileName)
	return os.Create(fileName)
}

func FileReadAll(path string) ([]byte, error) {
	cfg, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	b, e := ioutil.ReadAll(cfg)
	cfg.Close()
	return b, e
}

func FileScanner(path string) (*bufio.Scanner, *os.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return bufio.NewScanner(file), file, nil
}

func RemoveDirctory(path string) error {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	PthSep := string(os.PathSeparator)
	for _, info := range dirs {
		fpath := path + PthSep + info.Name()
		if info.IsDir() {
			err = RemoveDirctory(fpath)
			if err != nil {
				return err
			}
		} else {
			err = os.Remove(fpath)
			if err != nil {
				return err
			}
		}

	}
	return os.Remove(path)
}
