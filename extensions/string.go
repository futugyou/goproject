package extensions

import (
	"strconv"
	"strings"
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

func MaskString(s string, minLength int, maskRatio float64) string {
	length := len(s)
	if length < minLength {
		return strings.Repeat("*", length)
	}

	maskLength := int(float64(length) * maskRatio)
	startIndex := (length - maskLength) / 2

	maskedPart := strings.Repeat("*", maskLength)
	return s[:startIndex] + maskedPart + s[startIndex+maskLength:]
}

func StringToBoolPtr(s string) *bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return nil
	}
	return &b
}
