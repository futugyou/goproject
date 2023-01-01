package common

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 最大公因数，辗转相除法
func Gcd(a, b int) int {
	if a < b {
		return Gcd(b, a)
	}
	if b == 0 {
		return a
	}
	return Gcd(b, b%a)
}

// 最小公倍数
func Lcm(a, b int) int {
	return a * b / Gcd(a, b)
}

func Reverse(nums []int) {
	for l, r := 0, len(nums)-1; l < r; l, r = l+1, r-1 {
		nums[l], nums[r] = nums[r], nums[l]
	}
}
