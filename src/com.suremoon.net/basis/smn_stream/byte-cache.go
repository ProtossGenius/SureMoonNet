package smn_stream

import (
	"time"
	"sync"
	"sync/atomic"
)

type ByteCache struct {
	TimeOut time.Duration
	recvChan chan []byte
	last     []byte
	readLock sync.Mutex
	len   int32
}

func NewByteCache(writeTime int, timeout time.Duration) *ByteCache {
	return &ByteCache{TimeOut: timeout, recvChan:make(chan []byte, writeTime)}
}

func minInt(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func (this *ByteCache) Write(msg []byte)  {
	this.recvChan<-msg
	atomic.AddInt32(&this.len, int32(len(msg)))
}
//copy `size` byte from fConn.last to b[start:]
func (this *ByteCache) readCopy(b []byte, start, size int) (copyLen int) {
	copyLen = minInt(size, len(this.last))
	pos := start
	for i := 0; i < copyLen; i++ {
		b[pos] = this.last[i]
		pos++
	}
	this.last = this.last[copyLen:]
	return
}
func (this *ByteCache) Read(b []byte) (n int, err error) {
	this.readLock.Lock()
	defer this.readLock.Unlock()
	defer func() {	atomic.AddInt32(&this.len, int32(-n))
	}()
	size := len(b)
	acLen := 0
	if this.TimeOut == 0{
		this.TimeOut = 9 * time.Second
	}
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

func (this *ByteCache) Len() int {
	return int(this.len)
}