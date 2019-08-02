package smn_str

import "strings"

func DropLineComment(line string) string {
	line = strings.TrimSpace(line)
	return strings.Split(line, "//")[0]
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
