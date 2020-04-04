package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_err"
)

func check(err error) {
	smn_err.DftOnErr(err)
}

var pri *rsa.PrivateKey
var pub *rsa.PublicKey

func encode(bytes []byte) []byte {
	res, err := rsa.EncryptPKCS1v15(rand.Reader, pub, bytes)
	check(err)
	return res
}

func decode(bytes []byte) []byte {
	res, err := rsa.DecryptPKCS1v15(rand.Reader, pri, bytes)
	check(err)
	return res
}

var maxSize = map[int]int{
	2048: 245,
	4096: 501,
	1024: 117,
	512:  52,
}

func printTimeuse(pLen int) {
	var err error
	pri, err = rsa.GenerateKey(rand.Reader, pLen)
	check(err)
	pub = &pri.PublicKey
	start := time.Now().UnixNano()
	bLen := maxSize[pLen]
	bts := make([]byte, bLen)
	size := 1024 * 1024 / bLen
	oLen := 0
	for i := 0; i <= size; i++ {
		oLen = len(encode(bts))
	}
	end := time.Now().UnixNano()
	cost := end - start
	fmt.Println(pLen, "speed = ", 1024.0*1024.0/1000.0/(float64(cost)/1000000.0))
	fmt.Println("size change: ", float64(oLen)/float64(bLen))
}

func main() {
	printTimeuse(1024)
	printTimeuse(2048)
}
