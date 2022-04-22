package code0380

import (
	"fmt"
	"math/rand"
	"time"
)

func Exection() {
	r := RandomizedSet{
		nums:       make([]int, 0),
		valToIndex: make(map[int]int),
	}
	r.Insert(1)
	r.Insert(2)
	r.Insert(3)
	r.Insert(4)
	r.Insert(5)
	fmt.Println(r.GetRandom())
}

type RandomizedSet struct {
	nums       []int
	valToIndex map[int]int
}

func (r *RandomizedSet) Insert(val int) bool {
	if _, ok := r.valToIndex[val]; ok {
		return false
	}
	r.valToIndex[val] = len(r.nums)
	r.nums = append(r.nums, val)
	return true
}

func (r *RandomizedSet) Remove(val int) bool {
	if _, ok := r.valToIndex[val]; !ok {
		return false
	}
	index := r.valToIndex[val]
	last := r.nums[len(r.nums)-1]
	r.valToIndex[last] = index
	r.nums[len(r.nums)-1] = r.nums[index]
	r.nums[index] = last
	r.nums = r.nums[:len(r.nums)-1]
	delete(r.valToIndex, val)
	return true
}
func (r *RandomizedSet) GetRandom() int {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(r.nums)
	num := rr.Intn(n)
	//fmt.Println(num)
	return r.nums[num%n]
}
