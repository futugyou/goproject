package main

type itenant interface {
	paydeposit(deposit int)
	currmoney() int
}

type tenant struct {
	money int
}

func (t *tenant) paydeposit(deposit int) {
	t.money = t.money - deposit
}
func (t *tenant) currmoney() int {
	return t.money
}

type tenantA struct {
	tenant
}
type tenantB struct {
	tenant
}
