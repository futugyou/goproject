package extensions

import "strings"

func IndexWithOffset(s, substr string, offset int) int {
	if len(s) < offset {
		return -1
	}

	if idx := strings.Index(s[offset:], substr); idx >= 0 {
		return offset + idx
	}

	return -1
}
