package main

import (
	"com.suremoon.net/basis/smn_str_dict"
	"fmt"
)

func main() {
	sd := smn_str_dict.NewStrDictTree("abcd")
	sd.PutString("aa")
	sd.PutString("ab")
	sd.PutString("ac")
	sd.PutString("ad")
	for i := 0; i < 100; i++ {
		str := sd.RandStr(2)
		if str == "" {
			break
		}
		fmt.Println(i, "-----", str)
	}
}
