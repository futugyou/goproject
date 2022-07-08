package code0787

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

var memo [][]int
var indegree map[int][][]int
var srcc int

func Exection() {
	n := 3
	edges := [][]int{{0, 1, 100}, {1, 2, 100}, {0, 2, 500}}
	src := 0
	dst := 2
	k := 1
	r := exection(n, edges, src, dst, k)
	fmt.Println(r)
}

func exection(n int, edges [][]int, src, dst, k int) int {
	srcc = src
	k++
	memo = make([][]int, n)
	for i := 0; i < n; i++ {
		memo[i] = make([]int, k+1)
		for j := 0; j < k+1; j++ {
			memo[i][j] = -1
		}
	}
	indegree = map[int][][]int{}
	for _, v := range edges {
		from := v[0]
		to := v[1]
		price := v[2]
		indegree[to] = append(indegree[to], []int{from, price})
	}
	return dp(dst, k)
}

func dp(dst, k int) int {
	if dst == srcc {
		return 0
	}
	if k == 0 {
		return -1
	}
	if memo[dst][k] != -1 {
		return memo[dst][k]
	}
	res := math.MaxInt
	if curr, ok := indegree[dst]; ok {
		for _, v := range curr {
			from := v[0]
			price := v[1]
			subProblem := dp(from, k-1)
			if subProblem != -1 {
				res = common.Min(res, subProblem+price)
			}
		}
	}
	if res == math.MinInt {
		memo[dst][k] = -1
	} else {
		memo[dst][k] = res
	}
	return memo[dst][k]
}
