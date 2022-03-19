package code0875

import "fmt"

func Exection() {
	piles := []int{3, 6, 7, 11}
	h := 8
	exection(piles, h)
}

func exection(piles []int, h int) {
	if h < len(piles) {
		fmt.Println(-1)
		return
	}
	minspeed := 1
	maxspeed := max(piles)
	for {
		if minspeed > maxspeed {
			break
		}
		mid := minspeed + (maxspeed-minspeed)/2
		hour := execHour(piles, mid)
		if hour == h {
			maxspeed = mid - 1
		} else if hour > h {
			minspeed = mid + 1
		} else if hour < h {
			maxspeed = mid - 1
		}
	}
	fmt.Println(minspeed)
}

func execHour(piles []int, mid int) int {
	h := 0
	for i := 0; i < len(piles); i++ {
		if piles[i]%mid == 0 {
			h += piles[i] / mid
		} else {
			h += piles[i]/mid + 1
		}
	}
	return h
}

func max(piles []int) int {
	m := 0
	for i := 0; i < len(piles); i++ {
		if m < piles[i] {
			m = piles[i]
		}
	}
	return m
}
