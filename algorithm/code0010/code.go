package code0010

import "fmt"

func Exection() {
	s := "aa"
	p := "b"
	exection(s, p)
}

func exection(s, p string) {
	memo = map[string]bool{}
	r := dp(s, 0, p, 0)
	fmt.Println(r)
}

var memo map[string]bool

func dp(s string, i int, p string, j int) bool {
	m := len(s)
	n := len(p)
	// base case
	if j == n {
		return i == m
	}
	if i == m {
		if (n-j)%2 == 1 {
			return false
		}
		for ; j+1 < n; j += 2 {
			if p[j+1] != '*' {
				return false
			}
		}
		return true
	}

	key := fmt.Sprint(i, ",", j)
	if a, ok := memo[key]; ok {
		return a
	}

	res := false
	if s[i] == p[j] || p[j] == '.' {
		if j < n-1 && p[j+1] == '*' {
			res = dp(s, i, p, j+2) || dp(s, i+1, p, j)
		} else {
			res = dp(s, i+1, p, j+1)
		}
	} else {
		if j < n-1 && p[j+1] == '*' {
			res = dp(s, i, p, j+2)
		} else {
			res = false
		}
	}
	memo[key] = res
	return res

}
