package text

import "strings"

func AppendLine(builder *strings.Builder) *strings.Builder {
	builder.WriteByte('\n')
	return builder
}

func AppendStringLine(builder strings.Builder, value string) strings.Builder {
	builder.WriteString(value)
	builder.WriteByte('\n')
	return builder
}
