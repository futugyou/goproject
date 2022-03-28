package code0015

import (
	"fmt"
	"sort"
)

func Exection() {
	arr := []int{-1, 0, 1, 2, -1, -4}
	exection(arr)
}

func exection(arr []int) {
	sort.Ints(arr)
	fmt.Println(arr)
	result := nSumTarget(arr, 3, 0, 0)
	fmt.Println(result)
}

func nSumTarget(arr []int, n, start, target int) [][]int {
	result := make([][]int, 0)
	size := len(arr)
	if n < 2 || size < n {
		return result
	}
	if n == 2 {
		lo := start
		hi := size - 1
		for {
			if lo >= hi {
				break
			}
			left := arr[lo]
			right := arr[hi]
			sum := left + right
			if sum < target {
				for {
					if lo >= hi || arr[lo] != left {
						break
					}
					lo++
				}
			} else if sum > target {
				for {
					if lo >= hi || arr[hi] != right {
						break
					}
					hi--
				}
			} else {
				result = append(result, []int{left, right})
				for {
					if lo >= hi || arr[lo] != left {
						break
					}
					lo++
				}
				for {
					if lo >= hi || arr[hi] != right {
						break
					}
					hi--
				}
			}
		}
	} else {
		for i := start; i < size; i++ {
			sub := nSumTarget(arr, n-1, i+1, target-arr[i])
			for _, v := range sub {
				v = append(v, arr[i])
				result = append(result, v)
			}
			for {
				if i >= size-1 || arr[i] != arr[i+1] {
					break
				}
				i++
			}
		}
	}
	return result
}
