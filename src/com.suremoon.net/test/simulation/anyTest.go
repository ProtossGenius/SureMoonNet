package main

import (
	"com.suremoon.net/basis/smn_stream"
	"time"
	"fmt"
)

func main() {
	cache :=  smn_stream.NewByteCache(100, 10 * time.Second)
	str :="cache :=  smn_stream.NewByteCache(100, time.Second)"
	bstr := []byte(str)
	bts := make([]byte, len(bstr))
	go func() {
		time.Sleep(500 * time.Millisecond)
		cache.Write(bstr)
	}()
	cache.Read(bts)
	fmt.Println(string(bts))
}
