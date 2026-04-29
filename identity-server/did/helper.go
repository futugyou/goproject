package did

import (
	"encoding/hex"
	"strings"
)

func ToHex(value []byte, prefix bool) string {
	res := hex.EncodeToString(value)
	if prefix {
		return "0x" + res
	}
	return res
}

func HexToByteArray(value string) ([]byte, error) {
	value = strings.TrimPrefix(value, "0x")

	if len(value)%2 != 0 {
		value = "0" + value
	}

	return hex.DecodeString(value)
}
