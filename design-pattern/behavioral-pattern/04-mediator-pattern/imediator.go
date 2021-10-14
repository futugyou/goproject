package main

type imediator interface {
	addilandlord(landlord ilandlord)
	additenant(tenant itenant)
	renting()
}

type mediator struct {
	ilandlords []ilandlord
	itenants   []itenant
}

func (m *mediator) addilandlord(landlord ilandlord) {
	m.ilandlords = append(m.ilandlords, landlord)
}

func (m *mediator) additenant(tenant itenant) {
	m.itenants = append(m.itenants, tenant)
}

func (m *mediator) renting() {
	for _, landlord := range m.ilandlords {
		for _, tenant := range m.itenants {
			if landlord.rent() > tenant.currmoney() {
				tenant.paydeposit(landlord.rent())
			}
		}
	}
}
