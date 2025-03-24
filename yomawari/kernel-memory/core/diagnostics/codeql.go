package diagnostics

import (
	"strings"
)

func NLF(text string) string {
	if text == "" {
		return text
	}

	text = strings.ReplaceAll(text, "\n", "[char(10)]")
	text = strings.ReplaceAll(text, "\r", "[char(13)]")

	return text
}
