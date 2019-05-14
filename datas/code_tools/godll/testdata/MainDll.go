package main

// it is product by smnet.suremoon.com

import "C"
import pgt_factory "./hello/pgt_factory"

var valHello = pgt_factory.ProductHello()

//export helloa
func helloa(aaa string) string{
    return valHello.helloa (aaa)
}

//export hellob
func hellob(aaa, bbb string) string{
    return valHello.hellob (aaa, bbb)
}

//export helloc
func helloc(ls ...int){
    valHello.helloc (ls...)
}

