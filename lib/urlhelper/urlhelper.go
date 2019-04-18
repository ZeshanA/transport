package urlhelper

import (
	"strings"
)

func BuildQueryString(params map[string]string) string {
	count := len(params)
	if count == 0 {
		return ""
	}
	var sb = strings.Builder{}
	sb.WriteString("?")
	for key, value := range params {
		count -= 1
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		if count > 0 {
			sb.WriteString("&")
		}
	}
	return sb.String()
}
