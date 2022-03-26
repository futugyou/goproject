package code0986

import "fmt"

func Exection() {
	first := [][]int{{0, 2}, {5, 10}, {13, 23}, {24, 25}}
	second := [][]int{{1, 5}, {8, 12}, {15, 24}, {25, 26}}
	exection(first, second)
}

func exection(first [][]int, second [][]int) {
	len1 := len(first)
	len2 := len(second)
	a := 0
	b := 0
	result := make([][]int, 0)
	for {
		if a >= len1 || b >= len2 {
			break
		}
		if first[a][0] <= second[b][1] && first[a][1] >= second[b][0] {
			num1 := max(first[a][0], second[b][0])
			num2 := min(first[a][1], second[b][1])
			result = append(result, []int{num1, num2})
		}
		if first[a][1] < second[b][1] {
			a++
		} else {
			b++
		}
	}
	fmt.Println(result)
}

func min(i1, i2 int) int {
	if i1 > i2 {
		return i2
	}
	return i1
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}
