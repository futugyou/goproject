package code0300

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	arr := []int{10, 9, 2, 5, 3, 7, 101, 18}
	exection(arr)
}

func exection(arr []int) {
	n := len(arr)
	dp := make([]int, n)
	for i := 0; i < n; i++ {
		dp[i] = 1
	}
	for i := 0; i < n; i++ {
		for j := 0; j < i; j++ {
			if arr[i] > arr[j] {
				dp[i] = common.Max(dp[i], dp[j]+1)
			}
		}
	}
	res := 0
	for i := 0; i < n; i++ {
		res = common.Max(res, dp[i])
	}
	fmt.Println(res)
}
