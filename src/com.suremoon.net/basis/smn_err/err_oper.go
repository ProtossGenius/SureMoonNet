package smn_err

import "fmt"

func iserr(err error) bool {
	return err != nil
}

func noerr(err error) bool {
	return err == nil
}

type OnErr func(err error)

func DftOnErr(err error) {
	fmt.Println(err)
}