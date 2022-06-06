package code0752

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	deadends := []string{"0201", "0101", "0102", "1212", "2002"}
	target := "0202"
	r := exection(deadends, target)
	fmt.Println(r)
}

func exection(deadends []string, target string) int {
	deads := common.NewHashSet()
	for _, v := range deadends {
		deads.Add(v)
	}
	visited := common.NewHashSet()
	queue := common.NewQueue()
	step := 0
	visited.Add("0000")
	queue.Push("0000")
	for {
		if queue.Empty() {
			break
		}
		n := queue.Len()
		for i := 0; i < n; i++ {
			curr := queue.Pop().(string)
			if deads.Contains(curr) {
				continue
			}
			if curr == target {
				return step
			}

			for j := 0; j < 4; j++ {
				up := upwords(curr, j)
				if !visited.Contains(up) {
					visited.Add(up)
					queue.Push(up)
				}
				down := downwords(curr, j)
				if !visited.Contains(down) {
					visited.Add(down)
					queue.Push(down)
				}
			}
		}
		step++
	}
	return -1
}

func downwords(curr string, j int) string {
	t := []byte(curr)
	if t[j] == 0 {
		t[j] = 9
	}
	return string(t)
}

func upwords(curr string, j int) string {
	t := []byte(curr)
	if t[j] == 9 {
		t[j] = 0
	}
	return string(t)
}
