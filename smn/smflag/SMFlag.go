package smflag

import (
	"flag"
	"reflect"
)

type SMFlagValueItf interface {
	Parse()
}

type SMFlag struct {
	Value   SMFlagValueItf
	CmdList []string
}

func Parse(value SMFlagValueIt) *SMFlag {
	res := &SMFlag{Value: value}
	tv := reflect.TypeOf(value)
	flag.Parse()
	return res
}
