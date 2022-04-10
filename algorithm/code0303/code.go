package code0303

import "fmt"

func Exection() {
	arr := []int{4, 3, 1, 6, 8, 0, 5}
	exection(arr)
}

func exection(arr []int) {
	n := len(arr)
	pre := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pre[i] = pre[i-1] + arr[i-1]
	}
	fmt.Println(pre)
	a := 2
	b := 5
	r := pre[b] - pre[a]
	fmt.Println(r)
}
