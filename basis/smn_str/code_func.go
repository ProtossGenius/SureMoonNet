package smn_str

import "strings"

func DropLineComment(line string) string {
	line = strings.Split(line, "//")[0]
	return strings.TrimSpace(line)
}

func NotNullSpaceSplit(inp string) []string {
	strs := strings.Split(inp, " ")
	res := make([]string, 0, len(strs))
	for _, str := range strs {
		if str != "" {
			res = append(res, str)
		}
	}
	return res
}

//let first letter upper.   hello ->Hello
func InitialsUpper(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// PkgUpper from and_bnana to AndBnana.
func PkgUpper(pkg string) string {
	list := strings.Split(pkg, "_")
	for i := range list {
		list[i] = InitialsUpper(list[i])
	}

	return strings.Join(list, "_")
}

// ProtoUseDeal drop `[]`, `*` and let int as int64(proto not have int).
func ProtoUseDeal(typ string) (isArray bool, nt string) {
	if typ == "[]byte" {
		typ = "bytes"
	}

	isArray = strings.Contains(typ, "[]")
	typ = strings.ReplaceAll(typ, "[]", "")
	typ = strings.ReplaceAll(typ, "*", "")
	nt = strings.TrimSpace(typ)

	if nt == "int" {
		nt = "int64"
	}

	if nt == "uint" {
		nt = "uint64"
	}

	return
}
