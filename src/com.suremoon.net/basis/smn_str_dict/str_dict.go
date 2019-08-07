package smn_str_dict

/**
* smn_str_dict:
* A string dictionary for generating fixed-length and non-repeating random strings.
* Author ProtossGenius
*
*Use example:

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

*/

import (
	"math/rand"
)

type StrDictNode struct {
	isLeaf bool
	sons   []*StrDictNode
	sonId  map[rune]int
	idRune map[int]rune
	deep   int
	filled bool
}

func (this *StrDictNode) PutString(str string) {
	if len(str) <= this.deep {
		return
	}
	id := this.sonId[[]rune(str)[this.deep]]
	if this.sons[id] == nil {
		this.newSonById(id)
	}
	this.sons[id].PutString(str)
}

func (this *StrDictNode) newSonById(id int) {
	this.sons[id] = newStrDictNode(this.sonId, this.idRune, this.deep+1)
}

func (this *StrDictNode) Filled(deep int) bool {
	if this.filled {
		return true
	}
	if this.deep >= deep {
		return true
	}
	for _, node := range this.sons {
		if node == nil {
			return false
		}
		if !node.Filled(deep) {
			return false
		}
	}
	this.filled = true
	return this.filled
}

func (this *StrDictNode) RandStr(strLen int) string {
	if this.Filled(strLen) {
		return ""
	}
	unInit := make([]int, 0)
	unFilled := make([]int, 0)
	for i, node := range this.sons {
		if node == nil {
			unInit = append(unInit, i)
		} else if !node.Filled(strLen) {
			unFilled = append(unFilled, i)
		}
	}
	i := 0
	if len(unInit) != 0 {
		i = unInit[rand.Intn(len(unInit))]
		this.newSonById(i)
	} else {
		i = unFilled[rand.Intn(len(unFilled))]
	}
	lastRune := []rune(this.sons[i].RandStr(strLen))
	strRune := []rune{this.idRune[i]}
	return string(append(strRune, lastRune...))
}

func (this *StrDictNode) newSon(r rune) {
	this.newSonById(this.sonId[r])
}

func newStrDictNode(sonId map[rune]int, idRune map[int]rune, deep int) *StrDictNode {
	res := &StrDictNode{}
	res.deep = deep
	res.sonId = sonId
	res.sons = make([]*StrDictNode, len(sonId))
	res.idRune = idRune
	return res
}

func NewStrDictTree(letterSet string) *StrDictNode {
	sMap := make(map[rune]int)
	idRune := make(map[int]rune)
	i := 0
	for _, r := range letterSet {
		if _, ok := sMap[r]; ok {
			continue
		}
		sMap[r] = i
		idRune[i] = r
		i++
	}
	return newStrDictNode(sMap, idRune, 0)
}
