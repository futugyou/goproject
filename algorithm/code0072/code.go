package code0072

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	word1 := "intention"
	word2 := "execution"
	exection(word1, word2)
}

func exection(word1, word2 string) {
	m := len(word1)
	n := len(word2)
	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
	}
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for i := 0; i <= n; i++ {
		dp[0][i] = i
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = common.Min(dp[i-1][j], common.Min(dp[i][j-1], dp[i-1][j-1])) + 1
			}
		}
	}
	fmt.Println(dp)
	fmt.Println(dp[m][n])
}
