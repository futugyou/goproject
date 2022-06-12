package code0213

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{2, 7, 9, 3, 1}
	exection(nums)
}

func exection(nums []int) {
	n := len(nums)
	if n == 1 {
		fmt.Println(nums[0])
		return
	}
	if n == 2 {
		fmt.Println(common.Max(nums[0], nums[1]))
		return
	}
	dp := make([]int, n-1)
	dp[0] = nums[0]
	dp[1] = common.Max(nums[0], nums[1])
	for i := 2; i < n-1; i++ {
		dp[i] = common.Max(dp[i-1], dp[i-2]+nums[i])
	}
	dp1 := make([]int, n-1)
	dp1[0] = nums[1]
	dp1[1] = common.Max(nums[1], nums[2])
	for i := 3; i < n; i++ {
		dp1[i-1] = common.Max(dp1[i-2], dp1[i-3]+nums[i])
	}
	fmt.Println(common.Max(dp[n-2], dp1[n-2]))
}
