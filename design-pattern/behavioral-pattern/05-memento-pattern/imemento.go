package main

type imemento interface {
	creatememento(obj interface{})
	getmemento() interface{}
}

type memento struct {
	obj interface{}
}

func (m *memento) creatememento(obj interface{}) {
	m.obj = obj
}

func (m *memento) getmemento() interface{} {
	return m.obj
}
