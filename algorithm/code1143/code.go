package code1143

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	text1 := "abcde"
	text2 := "ace"
	r := exection(text1, text2)
	fmt.Println(r)
}

func exection(s1, s2 string) int {
	m := len(s1)
	n := len(s2)
	dp := make([][]int, m+1)
	for i := 0; i < m+1; i++ {
		dp[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] > s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = common.Max(dp[i][j-1], dp[i-1][j])
			}
		}
	}
	return dp[m][n]
}
