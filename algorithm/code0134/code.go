package code0134

import "fmt"

func Exection() {
	gas := []int{1, 2, 3, 4, 5}
	cost := []int{3, 4, 5, 1, 2}
	r := exection(gas, cost)
	fmt.Println(r)
}

func exection(gas []int, cost []int) int {
	n := len(gas)
	sum, minSum := 0, 0
	start := 0
	for i := 0; i < n; i++ {
		sum += gas[i] - cost[i]
		if sum < minSum {
			start = i + 1
			minSum = sum
		}
	}
	if sum < 0 {
		return -1
	}
	if start == n {
		return 0
	}
	return start
}
