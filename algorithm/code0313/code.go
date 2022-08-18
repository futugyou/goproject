package code0313

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 30
	primes := []int{2, 3, 5}
	r := exection(n, primes)
	fmt.Println(r)
}

func exection(n int, primes []int) int {
	nums := make([]int, n+1)
	pq := common.NewPriorityQueue2(cmpASC)
	for i := 0; i < len(primes); i++ {
		pq.Push([]int{1, primes[i], 1})
	}

	p := 1
	for {
		if p > n {
			break
		}
		curr := pq.Pop().([]int)
		val := curr[0]
		prime := curr[1]
		index := curr[2]

		if val != nums[p-1] {
			nums[p] = val
			p++
		}
		next := []int{nums[index] * prime, prime, index + 1}
		pq.Push(next)
		fmt.Println(nums)
	}
	return nums[n]
}

func cmpASC(a, b interface{}) int {
	x := a.([]int)
	y := b.([]int)
	if x[0] > y[0] {
		return 1
	} else if x[0] < y[0] {
		return -1
	}
	return 0
}
