package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.String("n", "name", "no use")
	flag.Parse()
	fmt.Println(flag.Args())
	fmt.Println("vim-go")
	idx := strings.Index("hello=\"world\"", "=")
	fmt.Println(idx)
}
