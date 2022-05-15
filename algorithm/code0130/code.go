package code0130

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	board := make([][]string, 9)
	for i := 0; i < 9; i++ {
		board[i] = make([]string, 8)
	}
	exection(board)
}

func exection(board [][]string) {
	row := len(board)
	col := len(board[0])
	fmt.Println(row, col)
	uf := common.NewUnionFind(row*col + 1)
	dummy := row * col
	// connect first/last col to dummy
	for i := 0; i < row; i++ {
		if board[i][0] == "o" {
			uf.Union(dummy, i*col)
		}
		if board[i][col-1] == "o" {
			uf.Union(dummy, i*col+col-1)
		}
	}
	// connect first/last row to dummy
	for i := 0; i < col; i++ {
		if board[0][i] == "o" {
			uf.Union(dummy, i)
		}
		if board[row-1][i] == "o" {
			uf.Union(dummy, col*(row-1)+i)
		}
	}
	d := [][]int{{1, 0}, {0, 1}, {0, -1}, {-1, 0}}
	for i := 1; i < row-1; i++ {
		for j := 1; j < col-1; j++ {
			if board[i][j] == "o" {
				for k := 0; k < 4; k++ {
					x := i + d[k][0]
					y := j + d[k][1]
					if board[x][y] == "o" {
						uf.Union(x*col+y, i*col+j)
					}
				}
			}
		}
	}
	for i := 1; i < row-1; i++ {
		for j := 1; j < col-1; j++ {
			if !uf.Connected(dummy, i*col+j) {
				board[i][j] = "x"
			}
		}
	}
}
