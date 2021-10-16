package main

func main() {
	s := newsubject()
	o1 := &observerA{}
	s.add(o1)
	s.setvalue(2)
}
