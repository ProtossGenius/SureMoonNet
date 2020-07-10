package main

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	_, err := os.Stat("./anyTest.go/Makefile")
	fmt.Println(err)

	fmt.Println("start .. ")
}
