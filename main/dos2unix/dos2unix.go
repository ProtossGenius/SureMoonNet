package main

import (
	"flag"
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
	path := flag.String("path", "./", "directory path.")
	exp := flag.String("exp", "txt", "exp. split with ',' ")
	deep := flag.Bool("deep", false, "is deep")
	flag.Parse()
	expList := strings.Split(*exp, ",")
	smn_file.DeepTraversalDir(*path, func(p string, info os.FileInfo) smn_file.FileDoFuncResult {
		if info.IsDir() && !*deep {
			return smn_file.FILE_DO_FUNC_RESULT_NO_DEAL
		}
		for _, e := range expList {
			if strings.HasSuffix(info.Name(), e) {
				all, err := smn_file.FileReadAll(p)
				check(err)
				newAll := strings.Replace(string(all), "\r\n", "\n", -1)
				f, err := smn_file.CreateNewFile(p)
				check(err)
				f.WriteString(newAll)
			}
		}
		return smn_file.FILE_DO_FUNC_RESULT_DEFAULT
	})

}
