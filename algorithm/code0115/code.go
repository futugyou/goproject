package code0115

import (
	"fmt"
)

func Exection() {
	s := "babgbag"
	t := "bag"
	r := exection(s, t)
	fmt.Println(r)
}

func exection(s, t string) int {
	m := len(s)
	n := len(t)

	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
		dp[i][n] = 1
	}

	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if s[i] == t[j] {
				dp[i][j] = dp[i+1][j+1] + dp[i+1][j]
			} else {
				dp[i][j] = dp[i+1][j]
			}
			fmt.Println(i, j, dp)
		}
	}
	fmt.Println(dp)
	fmt.Println(dp[0][0])

	memo = make([][]int, m)
	for i := 0; i < m; i++ {
		memo[i] = make([]int, n)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			memo[i][j] = -1
		}
	}
	return dpp(s, 0, t, 0)
}

var memo [][]int

func dpp(s string, i int, t string, j int) int {
	// base case 1
	if j == len(t) {
		return 1
	}
	// base case 2
	if len(s)-i < len(t)-j {
		return 0
	}
	if memo[i][j] != -1 {
		return memo[i][j]
	}

	res := 0
	if s[i] == t[j] {
		res += dpp(s, i+1, t, j+1) + dpp(s, i+1, t, j)
	} else {
		res += dpp(s, i+1, t, j)
	}
	memo[i][j] = res
	fmt.Println(memo)
	return res
}
