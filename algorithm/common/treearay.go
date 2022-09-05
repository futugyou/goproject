package common

var tree []int

func lowbit(a int) int {
	return a & (-a)
}

func UpdateTreeArray(x, val, n int) {
	tree = make([]int, n)
	for i := x; i < n; i += lowbit(i) {
		tree[i] += val
	}
}

func SumTreeArray(x int) int {
	r := 0
	for i := x; i >= 1; i -= lowbit(i) {
		r += tree[i]
	}
	return r
}
