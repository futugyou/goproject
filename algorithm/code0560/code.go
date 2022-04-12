package code0560

import "fmt"

func Exection() {
	arr := []int{1, 1, 1}
	k := 2
	exection(arr, k)
}

func exection(arr []int, k int) {
	n := len(arr)
	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + arr[i-1]
	}
	r := 0
	for i := 1; i <= n; i++ {
		for j := 0; j < i; j++ {
			if preSum[i]-preSum[j] == k {
				r++
			}
		}
	}
	fmt.Println(r)

	dic := make(map[int]int)
	dic[0] = 1
	r1 := 0
	sum0 := 0
	for i := 0; i < n; i++ {
		sum0 += arr[i]
		sum1 := sum0 - k
		if val, ok := dic[sum1]; ok {
			r1 += val
		}
		dic[sum0] = dic[sum0] + 1
	}
	fmt.Println(r1)
}
