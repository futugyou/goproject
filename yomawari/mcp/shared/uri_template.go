package shared

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func FormatUri(uriTemplate string, args map[string]any) (string, error) {
	var builder strings.Builder
	pos := 0
	for {
		start := strings.IndexByte(uriTemplate[pos:], '{')
		if start < 0 {
			builder.WriteString(uriTemplate[pos:])
			break
		}
		start += pos
		builder.WriteString(uriTemplate[pos:start])
		end := strings.IndexByte(uriTemplate[start:], '}')
		if end < 0 {
			return "", fmt.Errorf("unmatched '{' in URI template: %s", uriTemplate)
		}
		end += start
		expr := uriTemplate[start+1 : end]
		expansion, err := expandExpression(expr, args)
		if err != nil {
			return "", err
		}
		builder.WriteString(expansion)
		pos = end + 1
	}
	return builder.String(), nil
}

func expandExpression(expr string, args map[string]any) (string, error) {
	if expr == "" {
		return "", nil
	}
	var (
		prefix   string
		sep      string
		named    bool
		allowRes bool
	)
	mod := expr[0]
	switch mod {
	case '+':
		sep, allowRes = ",", true
		expr = expr[1:]
	case '#':
		prefix, sep, allowRes = "#", ",", true
		expr = expr[1:]
	case '.':
		prefix, sep = ".", "."
		expr = expr[1:]
	case '/':
		prefix, sep = "/", "/"
		expr = expr[1:]
	case ';':
		prefix, sep, named = ";", ";", true
		expr = expr[1:]
	case '?':
		prefix, sep, named = "?", "&", true
		expr = expr[1:]
	case '&':
		prefix, sep, named = "&", "&", true
		expr = expr[1:]
	default:
		sep = ","
	}

	parts := strings.Split(expr, ",")
	var result []string
	for _, part := range parts {
		name := part
		explode := false
		prefixLen := -1

		if strings.HasSuffix(name, "*") {
			explode = true
			name = name[:len(name)-1]
		} else if idx := strings.Index(name, ":"); idx >= 0 {
			nameOnly := name[:idx]
			lengthStr := name[idx+1:]
			n, err := strconv.Atoi(lengthStr)
			if err != nil {
				return "", fmt.Errorf("invalid prefix length: %s", name)
			}
			prefixLen = n
			name = nameOnly
		}

		val, ok := args[name]
		if !ok || val == nil {
			continue
		}

		switch v := val.(type) {
		case []string:
			if explode {
				for _, item := range v {
					enc := encode(item, allowRes)
					if named {
						result = append(result, fmt.Sprintf("%s=%s", name, enc))
					} else {
						result = append(result, enc)
					}
				}
			} else {
				var items []string
				for _, item := range v {
					items = append(items, encode(item, allowRes))
				}
				joined := strings.Join(items, ",")
				if named {
					result = append(result, fmt.Sprintf("%s=%s", name, joined))
				} else {
					result = append(result, joined)
				}
			}
		case map[string]string:
			if explode {
				for k, val := range v {
					result = append(result, fmt.Sprintf("%s=%s", encode(k, allowRes), encode(val, allowRes)))
				}
			} else {
				var items []string
				for k, val := range v {
					items = append(items, fmt.Sprintf("%s,%s", encode(k, allowRes), encode(val, allowRes)))
				}
				if named {
					result = append(result, fmt.Sprintf("%s=%s", name, strings.Join(items, ",")))
				} else {
					result = append(result, strings.Join(items, ","))
				}
			}
		default:
			str := fmt.Sprintf("%v", v)
			if prefixLen >= 0 && prefixLen < len(str) {
				str = str[:prefixLen]
			}
			enc := encode(str, allowRes)
			if named {
				if enc != "" {
					result = append(result, fmt.Sprintf("%s=%s", name, enc))
				} else {
					result = append(result, name)
				}
			} else {
				result = append(result, enc)
			}
		}
	}

	if len(result) == 0 {
		return "", nil
	}
	return prefix + strings.Join(result, sep), nil
}

func encode(s string, allowReserved bool) string {
	if allowReserved {
		return escapeExceptReserved(s)
	}
	return url.QueryEscape(s)
}

func escapeExceptReserved(s string) string {
	// Reserved characters (RFC 3986): :/?#[]@!$&'()*+,;=
	// We don't encode them if allowReserved=true
	return strings.ReplaceAll(url.PathEscape(s), "%2F", "/")
}
