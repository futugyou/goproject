package common

import "math"

// Rabin-Karp 指纹字符串查找算法
func RabinKarp(txt, pat string) int {
	// 位数
	L := len(pat)
	// 进制（只考虑 ASCII 编码）
	R := 256
	// 取一个比较大的素数作为求模的除数
	var Q float64 = 1658598167
	// R^(L - 1) 的结果
	var RL float64 = 1
	for i := 1; i <= L-1; i++ {
		// 计算过程中不断求模，避免溢出
		RL = math.Mod(RL*(float64)(R), Q)
	}
	// 计算模式串的哈希值，时间 O(L)
	var patHash float64 = 0
	for i := 0; i < len(pat); i++ {
		patHash = math.Mod(((float64)(R)*patHash + (float64)(pat[i])), Q)
	}

	// 滑动窗口中子字符串的哈希值
	var windowHash float64 = 0

	// 滑动窗口代码框架，时间 O(N)
	left := 0
	right := 0
	for {
		if right >= len(txt) {
			break
		}
		// 扩大窗口，移入字符
		windowHash = math.Mod((math.Mod(((float64)(R)*windowHash), Q) + (float64)(txt[right])), Q)
		right++

		// 当子串的长度达到要求
		if right-left == L {
			// 根据哈希值判断是否匹配模式串
			if windowHash == patHash {
				// 当前窗口中的子串哈希值等于模式串的哈希值
				// 还需进一步确认窗口子串是否真的和模式串相同，避免哈希冲突
				if pat == txt[left:right] {
					return left
				}
			}
			// 缩小窗口，移出字符
			windowHash = math.Mod((windowHash - math.Mod(((float64)(txt[left])*RL), Q) + Q), Q)
			// X % Q == (X + Q) % Q 是一个模运算法则
			// 因为 windowHash - (txt[left] * RL) % Q 可能是负数
			// 所以额外再加一个 Q，保证 windowHash 不会是负数

			left++
		}
	}
	// 没有找到模式串
	return -1
}
