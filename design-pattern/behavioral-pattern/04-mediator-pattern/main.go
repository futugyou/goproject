package main

func main() {
	ta := &tenantA{
		tenant: tenant{money: 3000},
	}
	tb := &tenantB{
		tenant: tenant{money: 4000},
	}
	ia := &landlordA{}
	ib := &landlordB{}

	m := &mediator{ilandlords: []ilandlord{ia, ib}, itenants: []itenant{ta, tb}}
	m.renting()
}
