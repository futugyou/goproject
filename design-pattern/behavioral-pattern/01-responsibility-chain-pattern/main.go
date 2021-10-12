package main

func main() {
	w := &work{worklevel: 12}
	s := &secondworkflow{next: nil, execlevel: 11}
	f := &firstflow{next: nil, execlevel: 10}
	f.setnextflow(s)
	f.exec(*w)
}
