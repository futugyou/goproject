package common

// https://labuladong.gitee.io/algo/3/28/97/
type Kmp struct {
	// dp[j][c] = next
	// 0 <= j < M，代表当前的状态
	// 0 <= c < 256，代表遇到的字符（ASCII 码）
	// 0 <= next <= M，代表下一个状态
	dp  [][]int
	pat string
}

func NewKmp(pat string) *Kmp {
	m := len(pat)
	// dp[状态][字符]=下一个状态
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, 256)
	}

	// base case
	dp[0][pat[0]] = 1
	// x 影子状态 初始0， 所谓影子状态，就是和当前状态具有相同的前缀
	x := 0
	// 当前状态i从1开始
	for i := 1; i < m; i++ {
		for c := 0; c < 256; c++ {
			if int(pat[i]) == c {
				dp[i][c] = i + 1
			} else {
				dp[i][c] = dp[x][c]
			}
		}
		// 更新影子状态
		x = dp[x][pat[i]]
	}
	return &Kmp{
		pat: pat,
		dp:  dp,
	}
}

func (k *Kmp) Search(txt string) int {
	m := len(k.pat)
	n := len(txt)
	// pat初始状态0
	j := 0
	for i := 0; i < n; i++ {
		// 计算 pat 的下一个状态
		j = k.dp[j][txt[i]]
		// 如果达到终止态，返回匹配开头的索引
		if j == m {
			return i - m + 1
		}
	}
	return -1
}
