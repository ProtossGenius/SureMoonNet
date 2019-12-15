package smn_file

import (
	"bufio"
	"io/ioutil"
	"os"
)

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func IsFileExist(fileName string) bool {
	var exist = true
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		exist = false
	}
	return exist
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
