package smn_stream

import (
	"bufio"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_str_rendering"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type ConditionFunc func(inp []byte) error

type ReadPipelineItf interface {
	Capture() error
	RemainingSize() int64 //unknow return -1,
	Read(buff []byte) (int, error)
	ConditionRead(condition ConditionFunc) ([]byte, error) //
	ByteBreakRead(condition ...byte) ([]byte, error)       //when get byte, end read
}

type ReadPipeline struct {
	reader    io.Reader
	closer    io.Closer
	readEnd   bool
	onceLock  sync.Mutex
	BuffSize  int //the size when read from reader.
	CacheSize int //the size save in chan.
	readChan  chan byte
	ErrChan   chan error
	TimeOut   time.Duration
}

func NewReadPipeline(reader io.Reader, closer io.Closer) *ReadPipeline {
	return &ReadPipeline{reader: reader, closer: closer}
}

func NewFileReadPipeline(fname string) (*ReadPipeline, error) {
	if !smn_file.IsFileExist(fname) {
		return nil, fmt.Errorf(smn_str_rendering.ERR_FILE_NOT_FOUND, fname)
	}
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(f)
	return NewReadPipeline(reader, f), nil
}

func (this *ReadPipeline) Capture() error {
	this.onceLock.Lock()
	defer this.onceLock.Unlock()
	this.readEnd = false
	if this.TimeOut <= 0 {
		this.TimeOut = 5 * time.Microsecond
	}
	if this.BuffSize <= 0 {
		this.BuffSize = 1024
	}
	if this.CacheSize <= 0 {
		this.CacheSize = 2048
	}
	this.readChan = make(chan byte, this.CacheSize)
	this.ErrChan = make(chan error, 10)
	go func() {
		readBuff := make([]byte, this.BuffSize)
		flag := true
		for flag {
			size, err := this.reader.Read(readBuff)
			if err != nil {
				if err.Error() == "EOF" {
					flag = false
				} else {
					this.ErrChan <- err
					return
				}
			}
			for i := 0; i < size; i++ {
				this.readChan <- readBuff[i]
			}
		}
		this.closer.Close()
		this.reader = nil
		this.readEnd = true
	}()
	return nil
}

func (this *ReadPipeline) RemainingSize() int64 {
	if this.readEnd && len(this.readChan) == 0 {
		return 0
	}
	return -1
}

func (this *ReadPipeline) read() (b byte, err error) {
	if len(this.readChan) != 0 {
		return <-this.readChan, nil
	}
	if this.readEnd {
		return b, errors.New("EOF")
	}
	time.Sleep(this.TimeOut)
	if len(this.readChan) == 0 {
		return b, errors.New(ErrTimeOut)
	}
	return <-this.readChan, nil
}
func (this *ReadPipeline) Read(buff []byte) (size int, err error) {
	if len(this.ErrChan) != 0 {
		return 0, <-this.ErrChan
	}
	size = len(buff)
	if len(this.readChan) < size {
		size = len(this.readChan)
	}
	for i := 0; i < size; i++ {
		buff[i] = <-this.readChan
	}
	return
}

func (this *ReadPipeline) ConditionRead(condition ConditionFunc) (res []byte, err error) {
	res = make([]byte, 0, 33)
	for {
		b, e := this.read()
		if e != nil {
			if e.Error() == "EOF" {
				e = nil
			}
			return res, e
		}
		res = append(res, b)
		if condition(res) == nil {
			return
		}
	}
	return res, condition(res)
}

func (this *ReadPipeline) ByteBreakRead(condition ...byte) (res []byte, err error) {
	res = make([]byte, 0, 33)
	for this.RemainingSize() != 0 {
		b, e := this.read()
		for _, c := range condition {
			if b == c {
				return
			}
		}
		if e != nil {
			return res, e
		}
		res = append(res, b)
	}
	return
}
