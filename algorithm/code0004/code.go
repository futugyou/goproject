package code0004

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums1 := []int{1, 3}
	nums2 := []int{2}
	r := exection(nums1, nums2)
	fmt.Println(r)
}

func exection(nums1 []int, nums2 []int) float64 {
	n1, n2 := len(nums1), len(nums2)
	if n1+n2 == 0 {
		return -1
	}
	if (n1+n2)%2 == 0 {
		l := findkth(nums1, nums1, (n1+n2)/2)
		r := findkth(nums1, nums2, (n1+n2)/2+1)
		return float64(l+r) / 2
	}
	return float64(findkth(nums1, nums1, (n1+n2)/2+1))
}

func findkth(nums1, nums2 []int, k int) int {
	n1, n2 := len(nums1), len(nums2)
	if n1 > n2 {
		n1, n2 = n2, n1
		nums1, nums2 = nums2, nums1
	}
	if n1 == 0 {
		return nums2[k-1]
	}
	if k == 1 {
		return common.Min(nums1[0], nums2[0])
	}
	k1 := common.Min(k/2, n1)
	k2 := k - k1
	switch {
	case nums1[k1-1] < nums2[k2-1]:
		return findkth(nums1[k1:], nums2, k2)
	case nums1[k1-1] > nums2[k2-1]:
		return findkth(nums1, nums2[k2:], k1)
	default:
		return nums1[k1-1]
	}
}
