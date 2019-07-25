package smn_muti_write_cache

import (
	"bytes"
	"io"
	"strings"
)

type FileMutiWriteCacheItf interface {
	WriteHead(str string) (int, error)
	WriteHeadLine(str string) (int, error)
	WriteTail(str string) (int, error)
	WriteTailLine(str string) (int, error)
	Append(itf ...FileMutiWriteCacheItf)
	GetFather() FileMutiWriteCacheItf
	SetFather(itf FileMutiWriteCacheItf)
	Output(oup io.Writer) (int, error)
}

type FileMutiWriteCache struct {
	Head     *bytes.Buffer
	Tail     *bytes.Buffer
	father   FileMutiWriteCacheItf
	Contains []FileMutiWriteCacheItf
}

func (this *FileMutiWriteCache) WriteHead(str string) (int, error) {
	return this.Head.WriteString(str)
}

func (this *FileMutiWriteCache) WriteTail(str string) (int, error) {
	return this.Tail.WriteString(str)
}

func (this *FileMutiWriteCache) WriteHeadLine(str string) (int, error) {
	return this.WriteHead(str + "\n")
}

func (this *FileMutiWriteCache) WriteTailLine(str string) (int, error) {
	return this.WriteTail(str + "\n")
}

func (this *FileMutiWriteCache) Append(itfs ...FileMutiWriteCacheItf) {
	this.Contains = append(this.Contains, itfs...)
	for _, itf := range itfs {
		itf.SetFather(this)
	}
}

func (this *FileMutiWriteCache) Output(oup io.Writer) (int, error) {
	_, err := oup.Write(this.Head.Bytes())
	if err != nil {
		return -1, err
	}
	for _, itf := range this.Contains {
		_, err := itf.Output(oup)
		if err != nil {
			return -1, err
		}
	}
	return oup.Write(this.Tail.Bytes())
}

func (this *FileMutiWriteCache) GetFather() FileMutiWriteCacheItf {
	return this.father
}

func (this *FileMutiWriteCache) SetFather(itf FileMutiWriteCacheItf) {
	this.father = itf
}

func NewFileMutiWriteCache() FileMutiWriteCacheItf {
	return &FileMutiWriteCache{father: nil, Head: bytes.NewBuffer(nil), Tail: bytes.NewBuffer(nil), Contains: make([]FileMutiWriteCacheItf, 0)}
}

type StringCache struct {
	val    string
	father FileMutiWriteCacheItf
}

func (this *StringCache) WriteHead(str string) (int, error) {
	this.val = str + this.val
	return len(str), nil
}

func (this *StringCache) WriteHeadLine(str string) (int, error) {
	return this.WriteHead(str + "\n")
}

func (this *StringCache) WriteTail(str string) (int, error) {
	this.val += str
	return len(str), nil
}

func (this *StringCache) WriteTailLine(str string) (int, error) {
	return this.WriteTail(str + "\n")
}

func (this *StringCache) Append(itf ...FileMutiWriteCacheItf) {
	panic("StringCache can't append")
}

func (this *StringCache) GetFather() FileMutiWriteCacheItf {
	return this.father
}

func (this *StringCache) SetFather(itf FileMutiWriteCacheItf) {
	this.father = itf
}

func (this *StringCache) Output(oup io.Writer) (int, error) {
	return oup.Write([]byte(this.val))
}

func NewStrCache(strs ...string) *StringCache {
	return &StringCache{val: strings.Join(strs, "")}
}
