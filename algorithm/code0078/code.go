package code0078

import "fmt"

func Exection() {
	nums := []int{1, 2, 3}
	exection(nums)
	fmt.Println(result)
}

var result [][]int

func exection(nums []int) {
	result = make([][]int, 0)
	path := make([]int, 0)
	backtrack(nums, 0, path)
}

func backtrack(nums []int, start int, path []int) {
	t := make([]int, len(path))
	copy(t, path)
	result = append(result, t)
	for i := start; i < len(nums); i++ {
		path = append(path, nums[i])
		backtrack(nums, i+1, path)
		path = path[:len(path)-1]
	}
}
