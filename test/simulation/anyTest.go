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
	str := "github.com/ProtossGenius/SureMoonNet/hellowrold"
	str = strings.Replace(str, "github.com/ProtossGenius/SureMoonNet", "github.com/ProtossGenius/SureMoonNet", -1)
	fmt.Println(str)
	smn_file.DeepTraversalDir("./", func(path string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
		}
		fData, err := smn_file.FileReadAll(path)
		check(err)
		fStr := string(fData)
		fStr = strings.Replace(fStr, "github.com/ProtossGenius/SureMoonNet", "github.com/ProtossGenius/SureMoonNet", -1)
		f, err := smn_file.CreateNewFile(path)
		check(err)
		f.WriteString(fStr)
		f.Close()
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

}
