package code0518

import (
	"fmt"
)

func Exection() {
	amount := 5
	coins := []int{1, 2, 5}
	r := exection(amount, coins)
	fmt.Println(r)
}

func exection(amount int, coins []int) int {
	n := len(coins)
	dp := make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dp[i] = make([]int, amount+1)
	}
	for i := 0; i < n+1; i++ {
		dp[i][0] = 1
	}
	for i := 1; i < n+1; i++ {
		for j := 1; j < amount+1; j++ {
			if j >= coins[i-1] {
				dp[i][j] = dp[i-1][j] + dp[i][j-coins[i-1]]
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[n][amount]
}
