package main

import (
	"basis/smn_file"
	"basis/smn_stream"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var filters = []string{".go", ".java"}

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
func getFileLines(fname string) int {
	fr, err := smn_stream.NewFileReadPipeline(fname)
	checkerr(err)
	checkerr(fr.Capture())
	count := 0
	for fr.RemainingSize() != 0 {
		fr.ByteBreakRead('\n')
		count++
	}
	return count
}
func main() {
	count := 0
	smn_file.DeepTraversalDir("./", func(path string, info os.FileInfo) bool {
		if info.IsDir() {
			return true
		}
		for _, n := range filters {
			if strings.HasSuffix(info.Name(), n) {
				count += getFileLines(path)
			}
		}
		return true
	})
	f, err := os.OpenFile("./code_line_statistics.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	checkerr(err)
	str := fmt.Sprintf("%s --code line statistics-- %d\n", time.Now().Format("2006-01-02 03:04:05.012 PM"), count)
	fmt.Println(str)
	io.WriteString(f, str)
}
