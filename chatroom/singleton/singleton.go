package singleton

type singleton struct {
	count int
}

var Instance = new(singleton)

func (s *singleton) Add() int {
	s.count++
	return s.count
}

// c:= singleton.Instance.Add()
