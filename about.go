package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	_, err := smn_file.DeepTraversalDir("./", func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}

		if !strings.HasSuffix(info.Name(), ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		fmt.Println(path)
		data, err := smn_file.FileReadAll(path)
		check(err)
		f, err := smn_file.CreateNewFile(path)
		check(err)
		defer f.Close()
		str := strings.ReplaceAll(string(data), "\"github.com/ProtossGenius/smnric/",
			`"github.com/ProtossGenius/smnric/`)
		str = strings.ReplaceAll(str, "snreader.", "snreader.")
		_, err = f.WriteString(str)
		check(err)

		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

	check(err)
}
