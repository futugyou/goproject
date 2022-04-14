package code1094

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	capacity := 4
	trips := [][]int{{2, 1, 5}, {3, 3, 7}}
	exection(trips, capacity)
}

func exection(trips [][]int, capacity int) {
	d := common.DiffArray{}
	arr := make([]int, 1000)
	d.Init(arr)
	for _, v := range trips {
		d.Increment(v[1], v[2]-1, v[0])
	}
	r := d.Result()
	fmt.Println(r)
	for _, v := range r {
		if v > capacity {
			fmt.Println(false)
			return
		}
	}
	fmt.Println(true)
}
