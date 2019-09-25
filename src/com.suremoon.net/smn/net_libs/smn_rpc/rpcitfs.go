package smn_rpc

import (
	"fmt"
	"net"
	"pb/base"
	"pb/dict"
	"github.com/golang/protobuf/proto"
	"com.suremoon.net/basis/smn_net"
)

func iserr(err error) bool {
	return err != nil
}
type ConnFunc func(conn net.Conn)

type MessageAdapterItf interface {
	WriteCall(dict dict.EDict, message proto.Message) (int, error)
	WriteRet(dict dict.EDict, message proto.Message, err error) (int, error)
	ReadCall() (*base.Call, error)
	ReadRet() (*base.Ret, error)
	GetConn() net.Conn
	Close() error
}

type MessageAdapter struct {
	c net.Conn
}

func NewMessageAdapter(conn net.Conn) MessageAdapterItf {
	return &MessageAdapter{c: conn}
}

func (this *MessageAdapter) Close() error {
	return this.c.Close()
}
func (this *MessageAdapter) GetConn() net.Conn {
	return this.c
}
func (this *MessageAdapter) WriteCall(dict dict.EDict, message proto.Message) (int, error) {
	bytes, err := proto.Marshal(message)
	if iserr(err) {
		return 0, err
	}
	msg := &base.Call{Dict: dict, Msg: bytes}
	bytes, err = proto.Marshal(msg)
	err = smn_net.WriteInt(len(bytes), this.c)
	if iserr(err) {
		return 0, err
	}
	return this.c.Write(bytes)
}

func (this *MessageAdapter) WriteRet(dict dict.EDict, message proto.Message, err error) (int, error) {
	bytes := make([]byte, 0)
	ret := &base.Ret{Dict: dict, Err: false}
	if err != nil {
		ret.Err = true
		bytes = []byte(err.Error())
	} else {
		var e error
		bytes, e = proto.Marshal(message)
		if e != nil {
			ret.Err = true
			bytes = []byte(e.Error())
		}
	}
	ret.Msg = bytes
	bytes, err = proto.Marshal(ret)
	err = smn_net.WriteInt(len(bytes), this.c)
	if iserr(err) {
		return 0, err
	}
	return this.c.Write(bytes)
}

func (this *MessageAdapter) ReadCall() (*base.Call, error) {
	len, err := smn_net.ReadInt(this.c)
	if iserr(err) {
		return nil, err
	}
	bytes := make([]byte, len)
	rl, err := this.c.Read(bytes)
	if err != nil {
		return nil, err
	}
	if rl != len {
		return nil, fmt.Errorf(smn_net.ErrNotGetEnoughLengthBytes, len, rl)
	}
	msg := &base.Call{}
	proto.Unmarshal(bytes, msg)
	return msg, err
}

func (this *MessageAdapter) ReadRet() (*base.Ret, error) {
	len, err := smn_net.ReadInt(this.c)
	if iserr(err) {
		return nil, err
	}
	bytes := make([]byte, len)
	rl, err := this.c.Read(bytes)
	if err != nil {
		return nil, err
	}
	if rl != len {
		return nil, fmt.Errorf(smn_net.ErrNotGetEnoughLengthBytes, len, rl)
	}
	msg := &base.Ret{}
	err = proto.Unmarshal(bytes, msg)
	return msg, err
}

func Int64ArrToIntArr(arr []int64) []int {
	res := make([]int, 0, len(arr))
	for _, i := range arr {
		res = append(res, int(i))
	}
	return res
}

func IntArrToInt64Arr(arr []int) []int64 {
	res := make([]int64, 0, len(arr))
	for _, i := range arr {
		res = append(res, int64(i))
	}
	return res
}

func UInt64ArrToUIntArr(arr []uint64) []uint {
	res := make([]uint, 0, len(arr))
	for _, i := range arr {
		res = append(res, uint(i))
	}
	return res
}

func UIntArrToUInt64Arr(arr []uint) []uint64 {
	res := make([]uint64, 0, len(arr))
	for _, i := range arr {
		res = append(res, uint64(i))
	}
	return res
}
