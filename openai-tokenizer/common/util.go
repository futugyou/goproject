package common

import (
	"strings"

	"github.com/dlclark/regexp2"
)

func IndexWithOffset(s, substr string, offset int) int {
	if len(s) < offset {
		return -1
	}

	if idx := strings.Index(s[offset:], substr); idx >= 0 {
		return offset + idx
	}

	return -1
}

func GetMapKeys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func IndexArrayWithOffset[T comparable](array []T, first T, offset int) int {
	if len(array) < offset {
		return -1
	}

	arr := array[offset:]
	for i := 0; i < len(arr); i++ {
		if arr[i] == first {
			return i + offset
		}
	}

	return -1
}

func ArrayFilter[T any](raws []T, filter func(T) bool) (ret []T) {
	for i := 0; i < len(raws); i++ {
		if filter(raws[i]) {
			ret = append(ret, raws[i])
		}
	}

	return
}

func Regexp2FindAllString(pattern, s string) []string {
	re := regexp2.MustCompile(pattern, 0)
	var matches []string
	m, _ := re.FindStringMatch(s)
	for m != nil {
		matches = append(matches, m.String())
		m, _ = re.FindNextMatch(m)
	}

	return matches
}
