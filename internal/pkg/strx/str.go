package strx

import (
	"regexp"
	"strings"
)

var chatRegexp = regexp.MustCompile("[.\\s]")

func FieldsX(str string) []string {
	var res []string
	for _, s := range chatRegexp.Split(str, -1) {
		s := strings.TrimSpace(s)
		if s != "" {
			res = append(res, s)
		}
	}
	return res
}
