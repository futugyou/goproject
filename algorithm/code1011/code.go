package code1011

import "fmt"

func Exection() {
	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	d := 5
	exection(weights, d)
}

func exection(weights []int, d int) {
	minhold := max(weights)
	maxhold := sum(weights)
	for {
		if minhold > maxhold {
			break
		}
		mid := minhold + (maxhold-minhold)/2
		day := execDay(weights, mid)
		if day == d {
			maxhold = mid - 1
		} else if day > d {
			minhold = mid + 1
		} else if day < d {
			maxhold = mid - 1
		}
	}
	fmt.Println(minhold)
}

func max(weights []int) int {
	m := 0
	for i := 0; i < len(weights); i++ {
		if m < weights[i] {
			m = weights[i]
		}
	}
	return m
}

func execDay(weights []int, mid int) int {
	d := 0
	for i := 0; i < len(weights); {
		t := mid
		for {
			if i >= len(weights) || t < weights[i] {
				break
			}
			t = t - weights[i]
			i++
		}
		d++
	}
	return d
}

func sum(piles []int) int {
	m := 0
	for i := 0; i < len(piles); i++ {
		m += piles[i]
	}
	return m
}
