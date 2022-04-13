package code0370

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	k := 5
	arr := [][]int{{1, 3, 2}, {2, 4, 3}, {0, 2, -2}}
	exection(arr, k)
}

func exection(arr [][]int, k int) {
	d := common.DiffArray{}
	init := make([]int, k)
	d.Init(init)
	for _, v := range arr {
		d.Increment(v[0], v[1], v[2])
	}
	fmt.Println(d.Result())
}
