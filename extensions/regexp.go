package extensions

import (
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
)

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

func Sanitize2String(input string, replace string) string {
	// Compile the regular expression to match non-alphanumeric characters
	re := regexp2.MustCompile(`[^a-zA-Z0-9]+`, regexp2.None)

	// Replace matches with a space
	result, _ := re.Replace(input, replace, 0, -1)

	// Trim and replace multiple spaces with a single space
	result = strings.TrimSpace(result)
	result = strings.Join(strings.Fields(result), replace)

	return result
}

func SanitizeString(input string, replace string) string {
	// Replace all non-alphanumeric characters with spaces
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	input = re.ReplaceAllString(input, replace)

	// Trim and replace multiple spaces with a single space
	input = strings.TrimSpace(input)
	input = strings.Join(strings.Fields(input), replace)

	return input
}
