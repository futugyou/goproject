package code0264

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 30
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	nums := make([]int, n+1)
	i2 := 1
	i3 := 1
	i5 := 1
	curr2 := 1
	curr3 := 1
	curr5 := 1
	p := 1
	for {
		if p > n {
			break
		}
		min := common.Min(curr2, common.Min(curr3, curr5))
		nums[p] = min
		p++
		if min == curr2 {
			curr2 = 2 * nums[i2]
			i2++
		}
		if min == curr3 {
			curr3 = 3 * nums[i3]
			i3++
		}
		if min == curr5 {
			curr5 = 5 * nums[i5]
			i5++
		}
		fmt.Println(nums)
	}
	return nums[n]
}
