package main

import (
	"basis/smn_directory"
	"basis/smn_stream"
	"fmt"
	"os"
	"strings"
)

var filters = []string{".go", ".java"}

func getFileLines(fname string) int {
	fr := smn_stream.FileReadPipeline{FileName: fname}
	fr.Capture()
	count := 0
	for fr.RemainingSize() != 0 {
		fr.ByteBreakRead('\n')
		count++
	}
	return count
}
func main() {
	count := 0
	smn_directory.DeepTraversalDir("./", func(path string, info os.FileInfo) bool {
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
	fmt.Println(count)
}
