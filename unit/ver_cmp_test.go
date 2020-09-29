package main

import (
	"testing"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_flag"
)

func assert(t *testing.T, chk bool) {
	if !chk {
		t.Fatal("false")
	}
}

func ver(v string) *smn_flag.Version {
	return new(smn_flag.Version).FromString(v)
}

func TestVerCmp(t *testing.T) {
	assert(t, ver("0.0.0.0").Less(ver("0.0.0.1")))
	assert(t, ver("0.0.0.1").Less(ver("0.0.1.0")))
	assert(t, ver("1.0.0.1").Less(ver("1.0.1.0")))
}
