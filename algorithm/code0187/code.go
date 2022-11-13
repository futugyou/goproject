package code0187

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

// Rabin-Karp https://labuladong.gitee.io/algo/2/20/28/
func Exection() {
	s := "AAAAACCCCCAAAAACCCCCCAAAAAGGGTTT"
	r := exection(s)
	fmt.Println(r)
}

func exection(s string) []string {
	n := len(s)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		if s[i] == 'A' {
			nums[i] = 0
		} else if s[i] == 'C' {
			nums[i] = 1
		} else if s[i] == 'G' {
			nums[i] = 2
		} else {
			nums[i] = 3
		}
	}
	seen := common.NewHashSet()
	res := make([]string, 0)
	l := 10
	r := 4
	rl := (int)(math.Pow((float64)(r), (float64)(l-1)))
	windowHash := 0
	left, right := 0, 0
	for {
		if right >= n {
			break
		}
		windowHash = r*windowHash + nums[right]
		right++
		if right-left == l {
			if seen.Contains(windowHash) {
				res = append(res, s[left:right-left])
			} else {
				seen.Add(windowHash)
			}
			windowHash = windowHash - nums[right]*rl
			left++
		}
	}
	return res
}
