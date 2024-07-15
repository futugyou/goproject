package extensions

import "github.com/dlclark/regexp2"

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
