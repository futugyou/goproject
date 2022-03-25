package code0042

import "fmt"

func Exection() {
	height := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	exection(height)
}

func exection(height []int) {
	n := len(height)
	result := 0
	leftmax := make([]int, n)
	rightmax := make([]int, n)
	leftmax[0] = height[0]
	rightmax[n-1] = height[n-1]
	for i := 1; i < n; i++ {
		leftmax[i] = max(height[i], leftmax[i-1])
	}
	for i := n - 2; i > 0; i-- {
		rightmax[i] = max(rightmax[i+1], height[i])
	}
	for i := 0; i < n; i++ {
		result += min(leftmax[i], rightmax[i]) - height[i]
	}
	fmt.Println(result)
}

func min(i1, i2 int) int {
	if i1 > i2 {
		return i2
	} else {
		return i1
	}
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}
