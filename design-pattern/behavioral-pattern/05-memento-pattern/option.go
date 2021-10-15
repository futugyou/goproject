package main

type option struct {
	age  int
	name string
	imementomanager
}

func (o *option) updateage(age int) {
	o.age = age
	memento := &memento{obj: *o}
	o.imementomanager.add(memento)
}
func (o *option) updatename(name string) {
	o.name = name
	memento := &memento{obj: *o}
	o.imementomanager.add(memento)
}

func (o *option) getlastchange() option {
	t := o.imementomanager.getlast()
	return t.getmemento().(option)
}
