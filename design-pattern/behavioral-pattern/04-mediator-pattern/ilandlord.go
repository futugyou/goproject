package main

type ilandlord interface {
	rent() int
}

type landlordA struct {
}

func (l *landlordA) rent() int {
	return 2000
}

type landlordB struct {
}

func (l *landlordB) rent() int {
	return 2100
}
