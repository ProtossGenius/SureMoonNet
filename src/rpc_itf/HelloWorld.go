package rpc_itf

import "pb/base"

type Login interface {
	DoLogin(user, pswd string, call base.Call) bool
}
