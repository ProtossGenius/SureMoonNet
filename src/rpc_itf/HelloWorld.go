package rpc_itf

type Login interface {
	DoLogin(user, pswd string, code int) (bool, int)
	//test array
	Test1(a []string, b []int, c []uint, d []uint64, e []int32) []int
}
