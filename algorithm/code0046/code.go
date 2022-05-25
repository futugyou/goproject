package code0046

import "fmt"

func Exection() {
	nums := []int{1, 2, 3}
	result = make([][]int, 0)
	exection(nums)
	fmt.Println(result)
}

var result [][]int

func exection(nums []int) {
	path := make([]int, 0)
	backtrack(nums, path)
}

func backtrack(nums []int, path []int) {
	if len(nums) == len(path) {
		result = append(result, path)
		return
	}
	for i := 0; i < len(nums); i++ {
		f := false
		for _, v := range path {
			if v == nums[i] {
				f = true
				break
			}
		}
		if f {
			continue
		}
		path = append(path, nums[i])
		backtrack(nums, path)
		path = path[:len(path)-1]
	}
}
