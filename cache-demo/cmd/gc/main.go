package main

import "time"

type Value struct {
	A string
	B int
	C time.Time
	D []byte
	E float32
	F *string
	T T
}
type T struct {
	H int
	I int
	K int
	L int
	M int
	N int
}

func main() {
	m := make(map[int]*Value, 10000000)
	for i := 0; i < 10000000; i++ {
		m[i] = &Value{}
	}

	for i := 0; ; i++ {
		delete(m, i)
		m[10000000+i] = &Value{}
		time.Sleep(6 * time.Millisecond)
	}
}
