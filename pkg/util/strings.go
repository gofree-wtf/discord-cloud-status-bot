package util

import "strings"

func SubstringAfter(str string, separator string) string {
	pos := strings.LastIndex(str, separator)
	if pos == -1 {
		return ""
	}

	adjustedPos := pos + len(separator)
	if adjustedPos >= len(str) {
		return ""
	}

	return str[adjustedPos:]
}
