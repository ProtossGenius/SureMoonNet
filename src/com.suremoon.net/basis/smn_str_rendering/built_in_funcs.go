package smn_str_rendering

import "sync"

type citf interface {
	Iadd(k string, v int) string
	Imult(k string, v int) string
	Idiv(k string, v int) string
	Isadd(k string, v int) int
	Ismult(k string, v int) int
	Isdiv(k string, v int) int
	Iget(k string) int
	Iset(k string, v int) string
	Inot0(k string) bool
}

type Counter struct {
	vs   map[string]int
	lock sync.Mutex
}

func (this *Counter) Isadd(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val += v
	this.vs[k] = val
	return val
}

func (this *Counter) Ismult(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val *= v
	this.vs[k] = val
	return v
}

func (this *Counter) Isdiv(k string, v int) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val /= v
	this.vs[k] = val
	return v
}

func (this *Counter) Iadd(k string, v int) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val += v
	this.vs[k] = val
	return ""
}

func (this *Counter) Imult(k string, v int) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val *= v
	this.vs[k] = val
	return ""
}

func (this *Counter) Idiv(k string, v int) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	val := this.vs[k]
	val /= v
	this.vs[k] = val
	return ""
}

func (this *Counter) Iget(k string) int {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.vs[k]
}

func (this *Counter) Inot0(k string) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.vs[k] != 0
}

func (this *Counter) Iset(k string, v int) string {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.vs[k] = v
	return ""
}
