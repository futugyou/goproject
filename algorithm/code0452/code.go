package code0452

import (
	"fmt"
	"sort"
)

func Exection() {
	points := [][]int{{10, 16}, {2, 8}, {1, 6}, {7, 12}}
	r := exection(points)
	fmt.Println(r)
}

func exection(points [][]int) int {
	sort.Slice(points, func(i, j int) bool {
		return points[i][1] < points[j][1]
	})
	count := 1
	end := points[0][1]
	for i := 1; i < len(points); i++ {
		curr := points[i]
		if curr[0] > end {
			count++
			end = curr[1]
		}
	}
	return count
}
