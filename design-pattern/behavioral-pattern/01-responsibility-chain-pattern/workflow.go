package main

type workflow interface {
	exec(w work)
	setnextflow(flow workflow)
}
