package code0295

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	m := NewMedianFinder()
	exection(m)
}

func exection(m MedianFinder) {
	m.AddNum(1)
	m.AddNum(2)
	//fmt.Println(m.FindMedian())
	m.AddNum(3)
	fmt.Println(m.FindMedian())
}

type MedianFinder struct {
	large *common.PriorityQueue
	small *common.PriorityQueue
}

func NewMedianFinder() MedianFinder {
	m := MedianFinder{}
	m.large = common.NewPriorityQueue(false)
	m.small = common.NewPriorityQueue(true)
	return m
}

func (m *MedianFinder) FindMedian() float32 {
	if m.large.Size() < m.small.Size() {
		t := m.small.Peek().(*common.PriorityQueueItem)
		return float32(t.Value)
	}
	if m.large.Size() > m.small.Size() {
		t := m.large.Peek().(*common.PriorityQueueItem)
		return float32(t.Value)
	}
	s := m.small.Peek().(*common.PriorityQueueItem)
	l := m.large.Peek().(*common.PriorityQueueItem)
	return (float32(s.Value) + float32(l.Value)) / 2
}
func (m *MedianFinder) AddNum(num int) {
	item := &common.PriorityQueueItem{
		Value:    num,
		Priority: num,
	}
	if m.small.Size() >= m.large.Size() {
		m.small.Push(item)
		m.large.Push(m.small.Pop())
	} else {
		m.large.Push(item)
		m.small.Push(m.large.Pop())
	}
}
