package code1109

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 5
	trips := [][]int{{1, 2, 10}, {2, 3, 20}, {2, 5, 25}}
	exection(trips, n)
}

func exection(trips [][]int, capacity int) {
	d := common.DiffArray{}
	arr := make([]int, capacity)
	d.Init(arr)
	for _, v := range trips {
		d.Increment(v[0]-1, v[1]-1, v[2])
	}
	r := d.Result()
	fmt.Println(r)
}
