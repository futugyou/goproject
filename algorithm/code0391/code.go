package code0391

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	rectangles := [][]int{{1, 1, 3, 3},
		{3, 1, 4, 2},
		{1, 3, 2, 4},
		{2, 2, 4, 4},
	}
	r := exection(rectangles)
	fmt.Println(r)
}

func exection(rectangles [][]int) bool {
	x1, y1 := math.MaxInt, math.MaxInt
	x2, y2 := -math.MaxInt, -math.MaxInt
	area := 0
	set := common.NewHashSet()
	for _, v := range rectangles {
		x1 = common.Min(x1, v[0])
		y1 = common.Min(y1, v[1])
		x2 = common.Max(x2, v[2])
		y2 = common.Max(y2, v[3])
		area += (x2 - x1) * (y2 - y1)
		if set.Contains([]int{x1, y1}) {
			set.Remove([]int{x1, y1})
		} else {
			set.Add([]int{x1, y1})
		}
		if set.Contains([]int{x1, y2}) {
			set.Remove([]int{x1, y2})
		} else {
			set.Add([]int{x1, y2})
		}
		if set.Contains([]int{x2, y1}) {
			set.Remove([]int{x2, y1})
		} else {
			set.Add([]int{x2, y1})
		}
		if set.Contains([]int{x2, y2}) {
			set.Remove([]int{x2, y2})
		} else {
			set.Add([]int{x2, y2})
		}
	}
	area1 := (x2 - x1) * (y2 - y1)
	if area != area1 || set.Size() != 4 ||
		set.Contains([]int{x1, y1}) ||
		set.Contains([]int{x2, y1}) ||
		set.Contains([]int{x1, y2}) ||
		set.Contains([]int{x2, y2}) {
		return false
	}
	return true
}
