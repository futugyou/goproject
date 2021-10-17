package main

func main() {
	con := statecontext{}
	state := &stateA{}
	state.exec(con)
}
