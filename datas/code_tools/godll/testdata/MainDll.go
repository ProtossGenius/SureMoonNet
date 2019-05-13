package main

// it is product by smnet.suremoon.com

import "C"
import pgt_factory "./hello/pgt_factory"

var valHello = pgt_factory.ProductHello()

func helloa(aaa string) string{
    return valHello.helloa (aaa)
}

func hellob(aaa, bbb string) string{
    return valHello.hellob (aaa, bbb)
}

func helloc(ls ...int){
    valHello.helloc (ls...)
}

