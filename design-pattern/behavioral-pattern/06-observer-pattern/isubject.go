package main

import "container/list"

type isubject interface {
	add(iobserver)
	remove(iobserver)
	notify()
}

type subject struct {
	observers *list.List
	value     int
}

func newsubject() *subject {
	s := new(subject)
	s.observers = list.New()
	return s
}

func (s *subject) getvalue() int {
	return s.value
}

func (s *subject) setvalue(v int) {
	s.value = v
	s.notify()
}

func (s *subject) add(observer iobserver) {
	s.observers.PushBack(observer)
}

func (s *subject) remove(observer iobserver) {
	for o := s.observers.Front(); o != nil; o = o.Next() {
		if o.Value.(*iobserver) == &observer {
			s.observers.Remove(o)
			break
		}
	}
}

func (s *subject) notify() {
	for o := s.observers.Front(); o != nil; o = o.Next() {
		o.Value.(iobserver).handler(s)
	}
}
