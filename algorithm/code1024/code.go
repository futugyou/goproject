package code1024

import (
	"fmt"
	"sort"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	clips := [][]int{{0, 2}, {4, 6}, {8, 10}, {1, 9}, {1, 5}, {5, 9}}
	t := 10
	r := exection(clips, t)
	fmt.Println(r)
}

func exection(clips [][]int, t int) int {
	sort.Slice(clips, func(i, j int) bool {
		if clips[i][0] == clips[j][0] {
			return clips[i][1] > clips[j][1]
		}

		return clips[i][0] < clips[j][0]
	})

	result := 1
	end := clips[0][1]
	tend := clips[0][1]
	for i := 1; i < len(clips); i++ {
		curr := clips[i]
		for {
			if !(curr[0] < end && end < curr[1]) {
				break
			}
			tend = common.Max(tend, curr[1])
		}
		end = tend
		result++
		if end >= t {
			return result
		}
	}
	return -1
}
