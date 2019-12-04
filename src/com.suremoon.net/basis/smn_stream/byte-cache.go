package smn_stream

import (
	"time"
	"sync"
	"sync/atomic"
	"github.com/pkg/errors"
)

const ErrByteCacheClosed = "ErrByteCacheClosed"

type iByteCache interface {
	Write(msg []byte)
	Read(b []byte) (n int, err error)
	Len() int
	SetTimeOut(t time.Duration)
}

func NewByteCache(writeTime int, timeout time.Duration) *ByteCache {
	return &ByteCache{c: &byteCacheWork{TimeOut: timeout, recvChan: make(chan []byte, writeTime)}}
}

type ByteCache struct {
	c iByteCache
}

func (this *ByteCache) Write(msg []byte) {
	this.c.Write(msg)
}

func (this *ByteCache) Read(b []byte) (n int, err error) {
	return this.c.Read(b)
}

func (this *ByteCache) Len() int {
	return this.c.Len()
}

func (this *ByteCache) Close() {
	this.c = bcc
}

func (this *ByteCache) ErrorClose(err error) {
	this.c = &byteCacheClose{err:err}
}

func (this *ByteCache) SetTimeOut(t time.Duration) {
	this.c.SetTimeOut(t)
}

var bcc = &byteCacheClose{err:errors.New(ErrByteCacheClosed)}

type byteCacheClose struct {
	err error
}

func (byteCacheClose) SetTimeOut(t time.Duration) {}

func (byteCacheClose) Write(msg []byte) {}

func (t *byteCacheClose) Read(b []byte) (n int, err error) {
	return 0, t.err
}

func (byteCacheClose) Len() int {
	return -1
}

type byteCacheWork struct {
	TimeOut  time.Duration
	recvChan chan []byte
	last     []byte
	readLock sync.Mutex
	len      int32
}

func (this *byteCacheWork) SetTimeOut(t time.Duration) {
	this.TimeOut = t
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (this *byteCacheWork) Write(msg []byte) {
	this.recvChan <- msg
	atomic.AddInt32(&this.len, int32(len(msg)))
}

//copy `size` byte from fConn.last to b[start:]
func (this *byteCacheWork) readCopy(b []byte, start, size int) (copyLen int) {
	copyLen = minInt(size, len(this.last))
	pos := start
	for i := 0; i < copyLen; i++ {
		b[pos] = this.last[i]
		pos++
	}
	this.last = this.last[copyLen:]
	return
}
func (this *byteCacheWork) Read(b []byte) (n int, err error) {
	this.readLock.Lock()
	defer this.readLock.Unlock()
	defer func() {
		atomic.AddInt32(&this.len, int32(-n))
	}()
	size := len(b)
	acLen := 0
	tout := time.Tick(this.TimeOut)
	for size > 0 {
		copyLen := this.readCopy(b, acLen, size)
		acLen += copyLen
		size -= copyLen
		if size == 0 {
			break
		}
		select {
		case this.last = <-this.recvChan:
		case <-tout:
			return acLen, nil
		}
	}
	return len(b), nil
}

func (this *byteCacheWork) Len() int {
	return int(this.len)
}
