package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

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

func NewCreateFile(fileName string) *os.File {
	RemoveFileIfExist(fileName)
	res, err := os.Create(fileName)
	checkerr(err)
	return res
}

func fileReadAll(path string) []byte {
	cfg, err := os.Open(path)
	checkerr(err)
	bytes, err := ioutil.ReadAll(cfg)
	checkerr(err)
	return bytes
}

type Replace struct {
	Key string
	Val string
}

func main() {
	flag.Parse()
	path := flag.Arg(0)
	if path == "" {
		path = "./初始化_腾讯_center.sql"
	}
	str := string(fileReadAll(path))
	cfg := string(fileReadAll("./init.cfg.txt"))
	for _, val := range strings.Split(cfg, "\n") {
		fmt.Println(val)
		val = strings.TrimSpace(val)
		if val == "" {
			continue
		}
		idx := strings.IndexByte(val, ':')
		str = strings.Replace(str, val[0:idx], val[idx+1:], -1)
	}
	f := NewCreateFile(path + ".out.sql")
	io.WriteString(f, str)
	//fmt.Println(str)
	f.Close()
	fmt.Println("press enter to exit")
	fmt.Scanln()
}
