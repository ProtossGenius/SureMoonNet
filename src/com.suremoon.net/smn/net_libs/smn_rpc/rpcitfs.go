package smn_rpc

import (
	"net"
	"pb/smn_base"
	"pb/smn_dict"
	"github.com/golang/protobuf/proto"
	"com.suremoon.net/basis/smn_net"
)

func iserr(err error) bool {
	return err != nil
}
type ConnFunc func(conn net.Conn)

type RpcSvrItf interface {
	OnMessage(c *smn_base.Call, conn net.Conn) (_d smn_dict.EDict, _p proto.Message, _e error)
}

type MessageAdapterItf interface {
	WriteCall(dict smn_dict.EDict, message proto.Message) (int, error)
	WriteRet(dict smn_dict.EDict, message proto.Message, err error) (int, error)
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

func (this *MessageAdapter) Close() error {
	return this.c.Close()
}
func (this *MessageAdapter) GetConn() net.Conn {
	return this.c
}

func WriteCall(conn net.Conn, dict smn_dict.EDict, message proto.Message) (int, error) {
	bts, err := proto.Marshal(message)
	if iserr(err) {
		return 0, err
	}
	msg := &smn_base.Call{Dict: dict, Msg: bts}
	bts, err = proto.Marshal(msg)
	return smn_net.WriteBytes(bts, conn)
}

func (this *MessageAdapter) WriteCall(dict smn_dict.EDict, message proto.Message) (int, error) {
	return WriteCall(this.c, dict, message)
}

func WriteRet(conn net.Conn, dict smn_dict.EDict, message proto.Message, err error) (int, error) {
	bts := make([]byte, 0)
	ret := &smn_base.Ret{Dict: dict, Err: false}
	if err != nil {
		ret.Err = true
		ret.Msg = []byte(err.Error())
	}else {
		ret.Msg, err = proto.Marshal(message)
		if err != nil {
			return 0, err
		}
	}
	bts, err = proto.Marshal(ret)
	if err != nil{
		return 0, err
	}
	return smn_net.WriteBytes(bts, conn)

}

func (this *MessageAdapter) WriteRet(dict smn_dict.EDict, message proto.Message, err error) (int, error) {
	return WriteRet(this.c, dict, message, err)
}

func (this *MessageAdapter) ReadCall() (*smn_base.Call, error) {
	bts, err := smn_net.ReadBytes(this.c)
	if err != nil{
		return nil, err
	}
	msg := &smn_base.Call{}
	proto.Unmarshal(bts, msg)
	return msg, err
}

func (this *MessageAdapter) ReadRet() (*smn_base.Ret, error) {
	bts, err := smn_net.ReadBytes(this.c)
	if err != nil{
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
