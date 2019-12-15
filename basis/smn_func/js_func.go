package smn_func

import (
	"github.com/robertkrimen/otto"
	"sync"
)

type jsFuncStruct struct {
	funcName string
	vm       *otto.Otto
}

type JsFunc func(params ...interface{}) (otto.Value, error)

func (this *jsFuncStruct) Call(params ...interface{}) (otto.Value, error) {
	return this.vm.Call(this.funcName, nil, params)
}

type JsFuncFactory struct {
	vm      *otto.Otto
	funcMap sync.Map
}

func NewJsFuncFactory(js string) (res *JsFuncFactory, err error) {
	res = &JsFuncFactory{}
	res.vm = otto.New()
	_, err = res.vm.Run(js)
	return res, err
}

func (this *JsFuncFactory) ProductFunc(funcName string) JsFunc {
	if val, ok := this.funcMap.Load(funcName); ok {
		return val.(JsFunc)
	}
	res := (&jsFuncStruct{vm: this.vm, funcName: funcName}).Call
	this.funcMap.Store(funcName, res)
	return res
}
