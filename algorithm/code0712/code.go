package code0712

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

var memo [][]int

func Exection() {
	s1 := "sea"
	s2 := "eat"
	memo = make([][]int, len(s1))
	for i := 0; i < len(s1); i++ {
		memo[i] = make([]int, len(s2))
		for j := 0; j < len(s2); j++ {
			memo[i][j] = -1
		}
	}
	r := exection(s1, 0, s2, 0)
	fmt.Println(r)
}

func exection(s1 string, i int, s2 string, j int) int {
	res := 0
	if i == len(s1) {
		for ; j < len(s2); j++ {
			res += (int)(s2[j])
		}
		return res
	}
	if j == len(s2) {
		for ; i < len(s1); i++ {
			res += (int)(s1[i])
		}
		return res
	}
	if memo[i][j] != -1 {
		return memo[i][j]
	}
	if s1[i] == s2[j] {
		memo[i][j] = exection(s1, i+1, s2, j+1)
	} else {
		memo[i][j] = common.Min((int)(s1[i])+exection(s1, i+1, s2, j), (int)(s2[j])+exection(s1, i, s2, j+1))
	}
	return memo[i][j]
}
