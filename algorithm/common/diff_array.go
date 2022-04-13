package common

type DiffArray struct {
	diff   []int
	result []int
}

func (d *DiffArray) Init(arr []int) {
	n := len(arr)
	d.diff = make([]int, n)
	d.result = make([]int, n)
	d.diff[0] = arr[0]
	for i := 1; i < n; i++ {
		d.diff[i] = arr[i] - arr[i-1]
	}
}

func (d *DiffArray) Increment(i, j, val int) {
	d.diff[i] += val
	if j+1 < len(d.diff) {
		d.diff[j+1] -= val
	}
}

func (d *DiffArray) Result() []int {
	d.result[0] = d.diff[0]
	for i := 1; i < len(d.diff); i++ {
		d.result[i] = d.diff[i] + d.result[i-1]
	}
	return d.result
}
