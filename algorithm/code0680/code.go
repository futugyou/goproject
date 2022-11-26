package code0680

import (
	"fmt"
)

func Exection() {
	s := "abc"
	r := validPalindrome(s)
	fmt.Println(r)
}

func validPalindrome(s string) bool {
	i := 0
	j := len(s) - 1

	if j == -1 || j == 0 || j == 1 {
		return true
	}

	for {
		if i >= j {
			return true
		}

		if s[i] != s[j] {
			return isPalindrome(s[i:j]) || isPalindrome(s[i+1:j+1])
		}
		i++
		j--
	}
}

func isPalindrome(s string) bool {
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
