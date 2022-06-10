package code0198

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	rooms := []int{2, 7, 9, 3, 1}
	exection(rooms)
}

func exection(rooms []int) {
	n := len(rooms)
	dp := make([]int, n)
	dp[0] = rooms[0]
	dp[1] = common.Max(rooms[1], rooms[0])
	for i := 2; i < n; i++ {
		dp[i] = common.Max(dp[i-1], dp[i-2]+rooms[i])
	}
	fmt.Println(dp[n-1])
}
