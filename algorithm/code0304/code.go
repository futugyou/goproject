package code0304

import "fmt"

func Exection() {
	arr := [][]int{{3, 0, 1, 4, 2}, {5, 6, 3, 2, 1}, {1, 2, 0, 1, 5}, {4, 1, 0, 1, 7}, {1, 0, 3, 0, 5}}
	exection(arr)
}

func exection(arr [][]int) {
	m := len(arr)
	n := len(arr[0])
	preSum := make([][]int, m+1)
	for i := 0; i < m+1; i++ {
		preSum[i] = make([]int, n+1)
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			preSum[i][j] = preSum[i-1][j] + preSum[i][j-1] - preSum[i-1][j-1] + arr[i-1][j-1]
		}
	}
	fmt.Println(preSum)
	// 2, 1, 4, 3
	x1 := 2
	y1 := 1
	x2 := 4
	y2 := 3
	r := preSum[x2+1][y2+1] - preSum[x1][y2+1] - preSum[x2+1][y1] + preSum[x1][y1]
	fmt.Println(r)
}
