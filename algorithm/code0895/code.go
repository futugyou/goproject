package code0895

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	f := &FreqStack{
		maxFreq:    0,
		valToFreq:  make(map[int]int),
		freqToVals: make(map[int]*common.Stack),
	}
	f.Push(1)
	f.Push(1)
	f.Push(1)
	f.Push(2)
	f.Push(2)
	f.Push(3)
	f.Push(4)
	fmt.Println(f.Pop())
	fmt.Println(f.Pop())
	fmt.Println(f.Pop())
	fmt.Println(f.Pop())
}

type FreqStack struct {
	maxFreq    int
	valToFreq  map[int]int
	freqToVals map[int]*common.Stack
}

func (f *FreqStack) Push(val int) {
	fr := f.valToFreq[val] + 1
	f.valToFreq[val] = fr
	if _, ok := f.freqToVals[fr]; !ok {
		f.freqToVals[fr] = common.NewStack()
	}
	f.freqToVals[fr].Push(val)
	if f.maxFreq < fr {
		f.maxFreq = fr
	}
}

func (f *FreqStack) Pop() int {
	vals := f.freqToVals[f.maxFreq]
	v := vals.Pop().(int)
	freq := f.valToFreq[v] - 1
	f.valToFreq[v] = freq
	if vals.Empty() {
		f.maxFreq--
	}
	return v
}
