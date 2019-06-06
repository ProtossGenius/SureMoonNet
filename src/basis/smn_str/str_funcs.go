package smn_str

import (
	"strings"
)

type CharReadEnd func(c rune) bool

var CIdentifierJoinEndCheck CharReadEnd = func(c rune) bool {
	return false
}
var CIdentifierDropEndCheck CharReadEnd = func(c rune) bool {
	if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_') || (c >= '0' && c <= '9') {
		return false
	}
	return true
}

func AnalysisTwoSplit(inp string, joinEnd, dropEnd CharReadEnd) (string, string) {
	idx := len(inp)
	for i, c := range inp {
		if joinEnd(c) {
			idx = i + 1
			break
		}
		if dropEnd(c) {
			idx = i
			break
		}
	}

	if idx > len(inp) {
		idx = len(inp) - 1
	}
	return inp[:idx], inp[idx:]
}

func AnalysisTwoSplitTrim(inp string, joinEnd, dropEnd CharReadEnd) (string, string) {
	inp = strings.TrimSpace(inp)
	a, b := AnalysisTwoSplit(inp, joinEnd, dropEnd)
	a = strings.TrimSpace(a)
	b = strings.TrimSpace(b)
	return a, b
}

func IsLowerChar(c rune) bool {
	return c >= 'a' && c <= 'z'
}

func IsUpperChar(c rune) bool {
	return c >= 'A' && c <= 'Z'
}

func IsChar(c rune) bool {
	return IsUpperChar(c) || IsLowerChar(c)
}

func CharUpper(c rune) rune {
	if IsLowerChar(c) {
		return c + 'A' - 'a'
	}
	return c
}
func CharLower(c rune) rune {
	if IsUpperChar(c) {
		return c + 'a' - 'A'
	}
	return c
}

func GetConstDefine(name string) string {
	res := ""
	for i, c := range name {
		if i == 0 {
			if !IsChar(c) {
				res += "_"
			} else {
				res += string(CharUpper(c))
			}
			continue
		}
		if IsLowerChar(c) {
			res += string(CharUpper(c))
		} else if IsUpperChar(c) {
			res += "_" + string(CharUpper(c))
		} else {
			if res[len(res)-1] != '_' {
				res += "_"
			}
		}
	}
	return res
}
