package ano_rpc_itf

import (
	"net"

	"github.com/ProtossGenius/SureMoonNet/pb/smn_base"
)

type Login interface {
	DoLogin(user, pswd string, code int) (bool, int)
	//test array
	Test1(a []string, b []int, c []uint, d []uint64, e []int32) []int
	Test2(key string, c net.Conn) bool
}

type MutiTest interface {
	SendMsg(mmm *smn_base.Call)
}
