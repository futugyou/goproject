package code0037

import "fmt"

func Exection() {
	board := make([][]int, 9)
	for i := 0; i < 9; i++ {
		board[i] = make([]int, 9)
	}
	exection(board)
}

var m, n = 9, 9

func exection(board [][]int) {
	r := backtrack(board, 0, 0)
	fmt.Println(r)
}

func backtrack(board [][]int, i, j int) bool {
	if j == n {
		return backtrack(board, i+1, 0)
	}
	if i == m {
		return true
	}
	if board[i][j] != -1 {
		return backtrack(board, i, j+1)
	}
	for c := i; c <= m; c++ {
		if !check(board, i, j, c) {
			continue
		}
		board[i][j] = c
		if backtrack(board, i, j+1) {
			return true
		}
		board[i][j] = -1
	}
	return false
}

func check(board [][]int, row, col, c int) bool {
	for i := 0; i < m; i++ {
		if board[i][col] == c {
			return false
		}
		if board[row][i] == c {
			return false
		}
		if board[(row/3)*3+i/3][(col/3)*3+i%3] == c {
			return false
		}
	}
	return true
}
