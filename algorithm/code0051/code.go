package code0051

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 4
	exection(n)
	fmt.Println(result)
}

var result [][][]string

func exection(n int) {
	result = make([][][]string, 0)
	board := make([][]string, n)
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
	}
	backtrack(board, 0)
}

func backtrack(board [][]string, row int) {
	n := len(board)
	if row == n {
		t := make([][]string, n)
		//copy(t, board)
		common.DeepCopy(&t, &board)
		result = append(result, t)
		return
	}
	for i := 0; i < n; i++ {
		if !check(board, row, i) {
			continue
		}
		board[row][i] = "*"
		backtrack(board, row+1)
		board[row][i] = ""
	}
}

func check(board [][]string, row, col int) bool {
	n := len(board)
	for i := 0; i < n; i++ {
		if board[i][col] == "*" {
			return false
		}
	}
	for i, j := row-1, col+1; i >= 0 && j < n; {
		if board[i][j] == "*" {
			return false
		}
		i--
		j++
	}
	for i, j := row-1, col-1; i >= 0 && j >= 0; {
		if board[i][j] == "*" {
			return false
		}
		i--
		j--
	}
	return true
}
