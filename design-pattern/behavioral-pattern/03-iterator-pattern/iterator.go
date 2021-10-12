package main

type iterator interface {
	hasnext() bool
	getnext() *user
}
