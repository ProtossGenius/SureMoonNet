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

func printTimeuse(pLen int) float64 {
	fmt.Println(pLen)
	var err error
	pri, err = rsa.GenerateKey(rand.Reader, pLen)
	check(err)
	pub = &pri.PublicKey
	start := time.Now().UnixNano()
	bLen := pLen/8 - 11
	bts := make([]byte, bLen)
	size := 1024 * 1024 / bLen
	oLen := 0
	for i := 0; i <= size; i++ {
		oLen = len(encode(bts))
	}
	end := time.Now().UnixNano()
	cost := end - start
	speed := 1024.0 * 1024.0 / 1000.0 / (float64(cost) / 1000000.0)
	fmt.Println(pLen, "speed = ", speed)
	fmt.Println("oLen: ", oLen)
	fmt.Println("size change: ", float64(oLen)/float64(bLen))
	return speed
}

func main() {
	/*
		var max float64
		mLen := 0
		for i := 1024; i < 2048; i++ {
			s := printTimeuse(i)
			if s > max {
				max, mLen = s, i
			}
		}
		fmt.Println(max, "   ", mLen)
	*/
	printTimeuse(1280) //max.
}
