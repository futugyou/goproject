package shared

import (
	"fmt"
	"regexp"
	"strings"
)

var exprPattern = regexp.MustCompile(`\{([+#./;?&]?)([A-Za-z0-9_.%]+(?:\*?)(?:,[A-Za-z0-9_.%]+\*?)*)\}`)

type UriParser struct {
	regex      *regexp.Regexp
	paramNames []string
}

func CreateUriParser(uriTemplate string) (*UriParser, error) {
	matches := exprPattern.FindAllStringSubmatchIndex(uriTemplate, -1)
	if matches == nil {
		return &UriParser{regex: regexp.MustCompile("^" + regexp.QuoteMeta(uriTemplate) + "$")}, nil
	}

	var pattern strings.Builder
	var paramNames []string = []string{}
	lastIndex := 0

	for _, match := range matches {
		start, end := match[0], match[1]
		opStart, opEnd := match[2], match[3]
		varsStart, varsEnd := match[4], match[5]

		pattern.WriteString(regexp.QuoteMeta(uriTemplate[lastIndex:start]))

		operator := uriTemplate[opStart:opEnd]
		vars := strings.Split(uriTemplate[varsStart:varsEnd], ",")

		switch operator {
		case "?": // Query
			pattern.WriteString(`(?:\?`)
			for i, v := range vars {
				if i > 0 {
					pattern.WriteString(`&?`)
				}
				paramNames = append(paramNames, v)
				pattern.WriteString(fmt.Sprintf(`(?:%s=([^/?&#]+))?`, regexp.QuoteMeta(v)))
			}
			pattern.WriteString(`)?`)
		case "&": // Query continuation
			pattern.WriteString(`(?:&`)
			for i, v := range vars {
				if i > 0 {
					pattern.WriteString(`&?`)
				}
				paramNames = append(paramNames, v)
				pattern.WriteString(fmt.Sprintf(`(?:%s=([^/?&#]+))?`, regexp.QuoteMeta(v)))
			}
			pattern.WriteString(`)?`)
		case "/": // Path
			for _, v := range vars {
				paramNames = append(paramNames, v)
				pattern.WriteString(`(?:/([^/?&#]+))?`)
			}
		case "#": // Fragment
			pattern.WriteString(`#`)
			for i, v := range vars {
				if i > 0 {
					pattern.WriteString(`,?`)
				}
				paramNames = append(paramNames, v)
				pattern.WriteString(`([^/?&#]+)?`)
			}
		case ".": // Label
			pattern.WriteString(`(?:\.`)
			for i, v := range vars {
				if i > 0 {
					pattern.WriteString(`\.?`)
				}
				paramNames = append(paramNames, v)
				pattern.WriteString(`([^/?&#.]+)?`)
			}
			pattern.WriteString(`)?`)
		case ";": // Path-style parameters
			for _, v := range vars {
				paramNames = append(paramNames, v)
				pattern.WriteString(fmt.Sprintf(`(?:;%s=([^/?&#;]+))?`, regexp.QuoteMeta(v)))
			}
		case "+", "": // Reserved or Simple string expansion
			for i, v := range vars {
				if i > 0 {
					pattern.WriteString(`,?`)
				}
				paramNames = append(paramNames, v)
				pattern.WriteString(`([^/?&#]+)?`)
			}
		default:
			return nil, fmt.Errorf("unsupported operator: %s", operator)
		}

		lastIndex = end
	}

	pattern.WriteString(regexp.QuoteMeta(uriTemplate[lastIndex:]))
	final := "^" + pattern.String() + "$"
	reg, err := regexp.Compile(final)
	if err != nil {
		return nil, err
	}
	return &UriParser{regex: reg, paramNames: paramNames}, nil
}

func (p *UriParser) Match(path string) map[string]string {
	m := p.regex.FindStringSubmatch(path)
	if m == nil {
		return nil
	}
	result := make(map[string]string)
	groupIndex := 1
	for _, name := range p.paramNames {
		if groupIndex < len(m) && m[groupIndex] != "" {
			result[name] = m[groupIndex]
		}
		groupIndex++
	}
	return result
}

func (p *UriParser) GetParamNames() []string {
	if p == nil {
		return []string{}
	}
	return p.paramNames
}
