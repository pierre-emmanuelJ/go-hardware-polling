package utils

import (
	"fmt"
	"strings"
)

func isSpace(c rune) bool {
	return (c == ' ' || c == '\t')
}

func StringRemoveAllWhiteSpace(s string) string {

	res := ""
	s = strings.TrimSpace(s)
	flag := false

	for _, c := range s {
		if isSpace(c) {
			flag = true
			continue
		}
		if flag {
			res = fmt.Sprintf("%s %c", res, c)
			flag = false
			continue
		}
		res = fmt.Sprintf("%s%c", res, c)
	}
	return res
}
