package word

import (
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ToUpper up
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower up
func ToLower(s string) string {
	return strings.ToLower(s)
}

// UnderscoreToUpperCamelCase UnderscoreToUpperCamelCase
func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "-", " ", -1)
	s = strings.Replace(s, "_", " ", -1)
	caser := cases.Title(language.English)
	s = caser.String(s)
	return strings.Replace(s, " ", "", -1)
}

// UnderscoreToLowerCamelCase UnderscoreToLowerCamelCase
func UnderscoreToLowerCamelCase(s string) string {
	s = UnderscoreToUpperCamelCase(s)
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// CamelCaseToUnderscore CamelCaseToUnderscore
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
