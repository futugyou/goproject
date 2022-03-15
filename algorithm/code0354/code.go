package code0354

import (
	"fmt"
	"sort"
)

func Exection() {
	envelopes := [][]int{{1, 8}, {2, 3}, {6, 4}, {5, 4}, {5, 2}, {6, 7}}
	exection(envelopes)
}

func exection(envelopes [][]int) {
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i][0] == envelopes[j][0] {
			return envelopes[i][1] > envelopes[j][1]
		} else {
			return envelopes[i][0] < envelopes[j][0]
		}
	})
	length := len(envelopes)
	arr := make([]int, length)
	for i := 0; i < length; i++ {
		arr[i] = envelopes[i][1]
	}
	fmt.Println(arr)
	dp := make([]int, length)
	for i := 0; i < length; i++ {
		dp[i] = 1
	}
	for i := 1; i < length; i++ {
		for j := 0; j < i; j++ {
			if arr[i] > arr[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
	}
	fmt.Println(dp)
}
func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
