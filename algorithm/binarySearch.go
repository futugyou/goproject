package main

import "fmt"

func binarySearch(list []int, target int) int {
	low := 0
	high := len(list) - 1

	step := 0
	for {
		step = step + 1
		if low <= high {
			mid := low + (high-low)/2
			guess := list[mid]
			if guess == target {
				fmt.Printf("count search %d \n", step)
				return mid
			}
			if guess > target {
				high = mid - 1
			} else {
				low = mid + 1
			}
		} else {
			return -1
		}
	}
}
