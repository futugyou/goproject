package code0066

import (
	"fmt"
)

func Exection() {
	digits := []int{9, 9, 9, 9}
	r := exection(digits)
	fmt.Println(r)
}

func exection(digits []int) []int {
	n := len(digits)
	if n == 0 {
		return digits
	}

	if digits[n-1] < 9 {
		digits[n-1] = digits[n-1] + 1
		return digits
	}

	for i := n - 2; i >= 0; i-- {
		digits[i+1] = 0
		if digits[i] < 9 {
			digits[i] = digits[i] + 1
			return digits
		}
		if i == 0 {
			digits[i] = 0
			return append([]int{1}, digits...)
		}
	}
	return digits
}
