package main

type statecontext struct {
	state istate
}

func (s *statecontext) set(state istate) {
	s.state = state
}

func (s *statecontext) doexec() {
	s.state.exec(*s)
}
