package code0315

import "fmt"

func Exection() {
	arr := []int{4, 3, 1, 6, 8, 0, 5}
	exection(arr)
}

type pair struct {
	val   int
	index int
}

var (
	temp   []pair
	result []int
)

func exection(arr []int) {
	n := len(arr)
	temp = make([]pair, n)
	result = make([]int, n)
	nums := make([]pair, n)
	for i := 0; i < n; i++ {
		nums[i] = pair{arr[i], i}
	}
	sort(nums, 0, n-1)
	fmt.Println(result)
}

func sort(nums []pair, lo, hi int) {
	if lo >= hi {
		return
	}
	mid := lo + (hi-lo)/2
	sort(nums, lo, mid)
	sort(nums, mid+1, hi)
	merge(nums, lo, mid, hi)
}

func merge(nums []pair, lo, mid, hi int) {
	for i := lo; i <= hi; i++ {
		temp[i] = nums[i]
	}
	i := lo
	j := mid + 1
	for x := lo; x <= hi; x++ {
		if i == mid+1 {
			nums[x] = temp[j]
			j++
		} else if j == hi+1 {
			nums[x] = temp[i]
			i++
			result[nums[x].index] += (j - mid - 1)
		} else if temp[i].val > temp[j].val {
			nums[x] = temp[j]
			j++
		} else {
			nums[x] = temp[i]
			i++
			result[nums[x].index] += (j - mid - 1)
		}
	}
}
