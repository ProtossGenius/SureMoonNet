package main

// it is product by smnet.suremoon.com

import "C"
import pgt_factory "../pgt_factory"
import pgt_interface "../pgt_interface"
import "sync"

var objarr = make([]pgt_interface.StrTagSysItf, 0, 30)
var idx int32 = 0
var oimap = make(map[int32]byte)
var oLock sync.Mutex

//export NewStrTagSysItf
func NewStrTagSysItf() int32 {
	oLock.Lock()
	defer oLock.Unlock()
	res := int32(0)
	if len(oimap) != 0 {
		for id := range oimap {
			res = id
			delete(oimap, id)
			break
		}
	} else {
		res = idx
		idx++
	}
	obj := pgt_factory.ProductStrTagSysItf()
	if res == int32(len(objarr)) {
		objarr = append(objarr, obj)
	} else {
		objarr[res] = obj
	}
	return res
}

//export DeleteStrTagSysItf
func DeleteStrTagSysItf(objid int32) bool {
	if objid >= int32(len(objarr)) || objarr[objid] == nil {
		return false
	}
	objarr[objid] = nil
	delete(oimap, objid)
	return true
}

//export AddTag
func AddTag(o_b_j_i_n_d_e_x int32, tagval, node string) (tagid int) {
	return objarr[o_b_j_i_n_d_e_x].AddTag(tagval, node)
}

//export PutTagStatus
func PutTagStatus(o_b_j_i_n_d_e_x int32, tagid, status int) (result bool) {
	return objarr[o_b_j_i_n_d_e_x].PutTagStatus(tagid, status)
}

//export PutTwoTagRelation
func PutTwoTagRelation(o_b_j_i_n_d_e_x int32, taga, tagb, relation int) {
	objarr[o_b_j_i_n_d_e_x].PutTwoTagRelation(taga, tagb, relation)
}

//export MayTowTagRealtion
func MayTowTagRealtion(o_b_j_i_n_d_e_x int32, taga, tagb, relation int) {
	objarr[o_b_j_i_n_d_e_x].MayTowTagRealtion(taga, tagb, relation)
}

//export GetTwoTagRelation
func GetTwoTagRelation(o_b_j_i_n_d_e_x int32, taga, tagb int) (relation int) {
	return objarr[o_b_j_i_n_d_e_x].GetTwoTagRelation(taga, tagb)
}

//export GetTagsetRelation
func GetTagsetRelation(o_b_j_i_n_d_e_x int32, tagid int) (relationJson string) {
	return objarr[o_b_j_i_n_d_e_x].GetTagsetRelation(tagid)
}

//export PutTagsetRelation
func PutTagsetRelation(o_b_j_i_n_d_e_x int32, tagid int, relationJson string) {
	objarr[o_b_j_i_n_d_e_x].PutTagsetRelation(tagid, relationJson)
}

//export GetFathers
func GetFathers(o_b_j_i_n_d_e_x int32, count, tagid int) (fatherid []int) {
	return objarr[o_b_j_i_n_d_e_x].GetFathers(count, tagid)
}

//export GetSons
func GetSons(o_b_j_i_n_d_e_x int32, count, tagid int) (sons []int) {
	return objarr[o_b_j_i_n_d_e_x].GetSons(count, tagid)
}

//export GetTagsets
func GetTagsets(o_b_j_i_n_d_e_x int32, num, id int) (sets []int) {
	return objarr[o_b_j_i_n_d_e_x].GetTagsets(num, id)
}

//export GetTagsetsFrom
func GetTagsetsFrom(o_b_j_i_n_d_e_x int32, Num, tagid int, fsets []int) (gsets []int) {
	return objarr[o_b_j_i_n_d_e_x].GetTagsetsFrom(Num, tagid, fsets)
}

func main() {
	// Need a main function to make CGO compile package as C shared library
}
