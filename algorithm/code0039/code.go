package code0039

import (
	"fmt"
	"sort"
)

func Exection() {
	candidates := []int{2, 3, 6, 7}
	target := 7
	sort.Ints(candidates)
	r := exection(candidates, target)
	fmt.Println(r)
}

func exection(candidates []int, target int) [][]int {
	result := make([][]int, 0)
	if len(candidates) == 0 {
		return result
	}

	sumTarget := target - candidates[0]
	switch {
	case sumTarget < 0:
		return result
	case sumTarget == 0:
		result = append(result, []int{candidates[0]})
	case sumTarget > 0:
		remains := exection(candidates, sumTarget)
		for _, v := range remains {
			way := append([]int{candidates[0]}, v...)
			result = append(result, way)
		}
	}

	result = append(result, exection(candidates[1:], target)...)
	return result
}
