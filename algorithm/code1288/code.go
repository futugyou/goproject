package code1288

import (
	"fmt"
	"sort"
)

func Exection() {
	intervals := [][]int{{1, 4}, {3, 6}, {2, 8}, {15, 18}}
	r := exection(intervals)
	fmt.Println(r)
}

func exection(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][0] == intervals[j][0] {
			return intervals[i][1] > intervals[j][1]
		} else {
			return intervals[i][0] < intervals[j][0]
		}
	})

	left := intervals[0][0]
	right := intervals[0][1]

	res := 0
	for i := 1; i < len(intervals); i++ {
		curr := intervals[i]
		if left <= curr[0] && right >= curr[1] {
			res++
		}
		if right >= curr[0] && right <= curr[1] {
			right = curr[1]
		}
		if right < curr[0] {
			left = curr[0]
			right = curr[1]
		}
	}
	return len(intervals) - res
}
