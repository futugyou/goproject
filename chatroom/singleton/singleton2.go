package singleton

import "sync"

type singleton2 struct {
	count int
}

var (
	instance2 *singleton2
	mutex     sync.Mutex
)

func New() *singleton2 {
	if instance2 == nil {
		mutex.Lock()
		if instance2 == nil {
			instance2 = new(singleton2)
		}
		mutex.Unlock()
	}
	return instance2
}
func (s *singleton2) Add() int {
	s.count++
	return s.count
}
