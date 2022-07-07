package code0887

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	k := 1
	n := 2
	r := exection(k, n)
	fmt.Println(r)
}

func exection(k, n int) int {
	dp := make([][]int, k+1)
	for i := 0; i < k+1; i++ {
		dp[i] = make([]int, n+1)
		for j := 0; j < n+1; j++ {
			dp[i][j] = math.MaxInt
		}
	}
	for i := 0; i < n+1; i++ {
		dp[0][i] = 0
		dp[1][i] = i
	}
	for i := 0; i < k+1; i++ {
		dp[i][0] = 0
	}
	for i := 2; i <= k; i++ {
		for j := 1; j <= n; j++ {
			for nn := 1; nn <= j; nn++ {
				dp[i][j] = common.Min(dp[i][j], common.Max(dp[i-1][nn-1], dp[i][j-nn]))
			}
		}
	}
	return dp[k][n]
}
