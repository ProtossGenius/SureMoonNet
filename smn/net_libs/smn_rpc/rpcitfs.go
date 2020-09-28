package smn_rpc

import (
	"net"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_net"
	"github.com/ProtossGenius/SureMoonNet/pb/smn_base"
	"google.golang.org/protobuf/proto"
)

func iserr(err error) bool {
	return err != nil
}

type ConnFunc func(conn net.Conn)

// StructCall  MessageAdapterItf's Param .
type StructCall struct {
	Dict int32
	Msg  proto.Message
}

// StructResult result .
type StructResult struct {
	Callback func(*smn_base.Ret)
	Ret      *smn_base.Ret
}

type RpcSvrItf interface {
	//@ret d -> dict, _p proto.Message, _e error
	OnMessage(c *smn_base.Call, conn net.Conn) (_d int32, _p proto.Message, _e error)
}

type MessageAdapterItf interface {
	WriteCall(dict int32, message proto.Message) (int, error)
	WriteRet(dict int32, message proto.Message, err error) (int, error)
	ReadCall() (*smn_base.Call, error)
	ReadRet() (*smn_base.Ret, error)
	GetConn() net.Conn
	Close() error
}

type MessageAdapter struct {
	c net.Conn
}

func NewMessageAdapter(conn net.Conn) MessageAdapterItf {
	return &MessageAdapter{c: conn}
}

func (ma *MessageAdapter) Close() error {
	return ma.c.Close()
}

func (ma *MessageAdapter) GetConn() net.Conn {
	return ma.c
}

func WriteCall(conn net.Conn, dict int32, message proto.Message) (int, error) {
	bts, err := proto.Marshal(message)
	if iserr(err) {
		return 0, err
	}

	msg := &smn_base.Call{Dict: dict, Msg: bts}
	bts, err = proto.Marshal(msg)

	if iserr(err) {
		return 0, err
	}

	return smn_net.WriteBytes(bts, conn)
}

func (ma *MessageAdapter) WriteCall(dict int32, message proto.Message) (int, error) {
	return WriteCall(ma.c, dict, message)
}

func WriteRet(conn net.Conn, dict int32, message proto.Message, err error) (int, error) {
	var bts []byte

	ret := &smn_base.Ret{Dict: dict, Err: false}

	if err != nil {
		ret.Err = true
		ret.Msg = []byte(err.Error())
	} else {
		ret.Msg, err = proto.Marshal(message)
		if err != nil {
			return 0, err
		}
	}

	bts, err = proto.Marshal(ret)
	if err != nil {
		return 0, err
	}

	return smn_net.WriteBytes(bts, conn)
}

func (ma *MessageAdapter) WriteRet(dict int32, message proto.Message, err error) (int, error) {
	return WriteRet(ma.c, dict, message, err)
}

func (ma *MessageAdapter) ReadCall() (*smn_base.Call, error) {
	bts, err := smn_net.ReadBytes(ma.c)
	if err != nil {
		return nil, err
	}

	msg := &smn_base.Call{}
	err = proto.Unmarshal(bts, msg)

	return msg, err
}

func (ma *MessageAdapter) ReadRet() (*smn_base.Ret, error) {
	bts, err := smn_net.ReadBytes(ma.c)
	if err != nil {
		return nil, err
	}

	msg := &smn_base.Ret{}
	err = proto.Unmarshal(bts, msg)

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
