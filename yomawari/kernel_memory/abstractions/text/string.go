package text

import "strings"

func NormalizeNewlines(text string, trim bool) string {
	if text == "" {
		return text
	}

	var builder strings.Builder
	builder.Grow(len(text))

	i := 0
	if trim {
		for i < len(text) && isWhitespace(text[i]) {
			i++
		}
	}

	lastNonWhitespacePos := -1

	for ; i < len(text); i++ {
		c := text[i]
		if c == '\r' {
			if i+1 < len(text) && text[i+1] == '\n' {
				i++
			}
			builder.WriteByte('\n')
		} else {
			builder.WriteByte(c)
		}

		if !trim || !isWhitespace(c) {
			lastNonWhitespacePos = builder.Len()
		}
	}

	if trim && lastNonWhitespacePos >= 0 {
		return builder.String()[:lastNonWhitespacePos]
	}

	return builder.String()
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func RemoveBOM(x string) string {
	x, _ = strings.CutPrefix(x, "\uFEFF")
	x, _ = strings.CutPrefix(x, "\u200B")
	return x
}
