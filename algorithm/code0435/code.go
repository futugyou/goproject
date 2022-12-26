package code0435

import (
	"fmt"
	"sort"
)

func Exection() {
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	r := exection(intervals)
	fmt.Println(r)
}

func exection(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][1] < intervals[j][1]
	})
	count := 1
	end := intervals[0][1]
	for i := 1; i < len(intervals); i++ {
		curr := intervals[i]
		if curr[0] >= end {
			count++
			end = curr[1]
		}
	}
	return count
}
