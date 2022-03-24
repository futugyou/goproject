package common

import (
	"container/list"
	"fmt"
)

type MonotonicQueue struct {
	arr *list.List
}

func (q *MonotonicQueue) New() {
	q.arr = list.New()
}

func (q *MonotonicQueue) Push(n int) {
	var prev *list.Element
	for e := q.arr.Back(); e != nil; e = prev {
		if e.Value.(int) < n {
			prev = e.Prev()
			q.arr.Remove(e)
		} else {
			break
		}
	}
	q.arr.PushBack(n)
	q.show()
}

func (q *MonotonicQueue) Max() int {
	if e := q.arr.Front(); e != nil {
		return e.Value.(int)
	}
	return -1
}
func (q *MonotonicQueue) Pop(n int) {
	if e := q.arr.Front(); e != nil {
		if e.Value.(int) == n {
			q.arr.Remove(e)
		}
	}
}

func (q *MonotonicQueue) show() {
	for e := q.arr.Front(); e != nil; e = e.Next() {
		fmt.Print(e.Value, " ")
	}
	fmt.Println()
}
