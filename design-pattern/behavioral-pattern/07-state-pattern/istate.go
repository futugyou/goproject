package main

type istate interface {
	exec(con statecontext)
}
