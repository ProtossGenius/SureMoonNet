package smn_flag

import (
	"fmt"
	"strconv"
	"strings"
)

// Version version.
type Version struct {
	FirstNo  int
	SecondNo int
	ThridNo  int
	ForthNo  int
}

// Less .
func (v *Version) Less(rhs *Version) bool {
	if v.FirstNo != rhs.FirstNo {
		return v.FirstNo < rhs.FirstNo
	}
	if v.SecondNo != rhs.SecondNo {
		return v.SecondNo < rhs.SecondNo
	}
	if v.ThridNo != rhs.ThridNo {
		return v.ThridNo < rhs.ThridNo
	}

	return v.ForthNo < rhs.ForthNo
}

// NewVersion create a Version.
func NewVersion(v1, v2, v3, v4 int) *Version {
	return &Version{v1, v2, v3, v4}
}

// ToString  v to string.
func (v *Version) ToString() string {
	return fmt.Sprintf("%d.%d.%d.%d", v.FirstNo, v.SecondNo, v.ThridNo, v.ForthNo)
}

func toint(num string) int {
	i, _ := strconv.Atoi(num)
	return i
}

// FromString init from string.
func (v *Version) FromString(ver string) *Version {
	list := strings.Split(ver, ".")
	for len(list) < 4 {
		list = append(list, "0")
	}

	v.FirstNo = toint(list[0])
	v.SecondNo = toint(list[1])
	v.ThridNo = toint(list[2])
	v.ForthNo = toint(list[3])

	return v
}

// RegisterVersion version flag.
func RegisterVersion(name string, v Version, infos ...string) {
	RegisterBool("v", "show version", func(val string) error {
		fmt.Println(name)
		for _, info := range infos {
			fmt.Println(info)
		}

		fmt.Println(v.ToString())

		return nil
	})
}
