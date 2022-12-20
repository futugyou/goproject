package code0056

import (
	"fmt"
	"sort"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	r := exection(intervals)
	fmt.Println(r)
}

func exection(intervals [][]int) [][]int {
	result := make([][]int, 0)
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})
	result = append(result, intervals[0])
	for i := 1; i < len(intervals); i++ {
		curr := intervals[i]
		last := result[len(result)-1]
		if curr[0] <= last[1] {
			result[len(result)-1][1] = common.Max(curr[1], last[1])
		} else {
			result = append(result, intervals[i])
		}
	}
	return result
}
