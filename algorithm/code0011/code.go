package code0011

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	height := []int{1, 8, 6, 2, 5, 4, 8, 3, 7}
	r := exection(height)
	fmt.Println(r)
}

func exection(height []int) int {
	count := 0
	left, right := 0, len(height)-1
	for left < right {
		w := right - left
		h := common.Min(height[left], height[right])
		count = common.Max(count, w*h)
		if height[left] <= height[right] {
			left++
		} else {
			right--
		}
	}
	return count
}
