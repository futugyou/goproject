package code0205

import "fmt"

func Exection() {
	s := "egg"
	t := "add"
	r := isIsomorphic(s, t)
	fmt.Println(r)
}

func isIsomorphic(s string, t string) bool {
	dic := make(map[byte]byte, len(s))

	for i := 0; i < len(s); i++ {
		if val, ok := dic[s[i]]; ok {
			if val != s[i]-t[i] {
				return false
			}
		} else {
			dic[s[i]] = s[i] - t[i]
		}
	}

	dic = make(map[byte]byte, len(s))
	for i := 0; i < len(s); i++ {
		if val, ok := dic[t[i]]; ok {
			if val != t[i]-s[i] {
				return false
			}
		} else {
			dic[t[i]] = t[i] - s[i]
		}
	}

	return true
}
