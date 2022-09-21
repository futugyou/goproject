package code0233

import (
	"fmt"
	"strconv"
)

func Exection() {
	n := 10000
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	s := strconv.Itoa(n)
	m := len(s)
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, m)
		for j := range dp[i] {
			dp[i][j] = -1
		}
	}
	var f func(int, int, bool) int
	f = func(i, cnt int, isLimit bool) (res int) {
		if i == m {
			return cnt
		}
		if !isLimit && dp[i][cnt] > 0 {
			return dp[i][cnt]
		}
		max := 9
		if isLimit {
			max = int(s[i] - '0')
		}
		for j := 0; j <= max; j++ {
			c := cnt
			if j == 1 {
				c++
			}
			res += f(i+1, c, isLimit && j == max)
		}

		if !isLimit {
			dp[i][cnt] = res
		}
		return
	}
	return f(0, 0, true)
}
