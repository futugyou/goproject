package code0239

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k := 3
	exection(nums, k)
}

func exection(nums []int, k int) {
	q := common.MonotonicQueue{}
	q.New()
	result := make([]int, 0)
	for i := 0; i < len(nums); i++ {
		if i < k-1 {
			q.Push(nums[i])
		} else {
			q.Push(nums[i])
			result = append(result, q.Max())
			q.Pop(nums[i-k+1])
		}
	}
	fmt.Println(result)
}
