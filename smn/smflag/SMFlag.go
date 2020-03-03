package smflag

import (
	"flag"
)

type SMFlagValueItf interface {
	Parse()
}

type SMFlag struct {
	Value   SMFlagValueItf
	CmdList []string
}

func Parse(value SMFlagValueItf) *SMFlag {
	res := &SMFlag{Value: value}
	// tv := reflect.TypeOf(value)
	flag.Parse()
	return res
}
