package code0322

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	coins := []int{1, 2, 5}
	amount := 11
	exection(coins, amount)
}

func exection(coins []int, amount int) {
	dp := make([]int, amount+1)
	for i := 0; i < len(dp); i++ {
		dp[i] = amount + 1
	}
	dp[0] = 0
	for i := 1; i < amount+1; i++ {
		for _, coin := range coins {
			if i >= coin {
				dp[i] = common.Min(dp[i], dp[i-coin]+1)
			}
		}
	}
	fmt.Println(dp)
	fmt.Println(dp[len(dp)-1])
}
