package common

import (
	"container/list"
)

type Queue struct {
	list *list.List
}

func NewQueue() *Queue {
	list := list.New()
	return &Queue{list}
}

func (queue *Queue) Push(value interface{}) {
	queue.list.PushBack(value)
}

func (queue *Queue) Pop() interface{} {
	e := queue.list.Front()
	if e != nil {
		queue.list.Remove(e)
		return e.Value
	}
	return nil
}

func (queue *Queue) Peak() interface{} {
	e := queue.list.Front()
	if e != nil {
		return e.Value
	}

	return nil
}

func (queue *Queue) Len() int {
	return queue.list.Len()
}

func (queue *Queue) Empty() bool {
	return queue.list.Len() == 0
}
