package main

// it is product by smnet.suremoon.com

import "C"
import pgt_factory "./hello/pgt_factory"
import pgt_interface "./hello/pgt_interface"
import "sync"
var objarr = make([]pgt_interface.Hello, 0, 30)
var idx int32 = 0
var oimap = make(map[int32]byte)
var oLock sync.Mutex

//export NewHello
func NewHello() int32{
    oLock.Lock()
    defer oLock.Unlock()
    res := int32(0)
    if len(oimap) != 0{
        for id := range oimap {
            res = id
            delete(oimap, id)
            break
        }
    }else {
        res = idx
        idx++
    }
    obj := pgt_factory.ProductHello()
    if res == int32(len(objarr)){
        objarr = append(objarr, obj)
    }else {
        objarr[res] = obj
    }
    return res
}

//export DeleteHello
func DeleteHello(objid int32) bool {
    if objid >= int32(len(objarr)) || objarr[objid] == nil{
        return false
    }
    objarr[objid] = nil
    delete(oimap, objid)
    return true
}



//export Helloelloa
func Helloelloa(o_b_j_i_n_d_e_x int32, aaa string) string{
    return objarr[o_b_j_i_n_d_e_x].Helloelloa ( aaa)
}

//export Helloellob
func Helloellob(o_b_j_i_n_d_e_x int32, aaa, bbb string) string{
    return objarr[o_b_j_i_n_d_e_x].Helloellob ( aaa, bbb)
}

//export Helloelloc
func Helloelloc(o_b_j_i_n_d_e_x int32, ls ...int){
    objarr[o_b_j_i_n_d_e_x].Helloelloc ( ls...)
}


func main() {
    // Need a main function to make CGO compile package as C shared library
}

