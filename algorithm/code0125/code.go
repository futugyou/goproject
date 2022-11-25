package code0125

import (
	"fmt"
	"strings"
)

func Exection() {
	s := "A man, a plan, a canal: Panama"
	r := isPalindrome(s)
	fmt.Println(r)
}

func isPalindrome(s string) bool {
	s = strings.ToLower(strip(s))
	i := 0
	j := len(s) - 1

	if j == 0 || j == -1 {
		return true
	}
	for {
		if s[i] != s[j] {
			return false
		}
		if i >= j {
			return true
		}
		i++
		j--
	}
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') {
			result.WriteByte(b)
		}
	}
	return result.String()
}
