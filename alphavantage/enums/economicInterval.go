package enums

type EconomicGdpInterval interface {
	privateeconomicGdpInterval()
	String() string
}

type economicGdpInterval string

func (c economicGdpInterval) privateeconomicGdpInterval() {}
func (c economicGdpInterval) String() string {
	return string(c)
}

const EconomicGdpQuarterly economicGdpInterval = "quarterly"
const EconomicGdpAnnual economicGdpInterval = "annual"

type EconomicTreasuryInterval interface {
	privateeconomicTreasuryInterval()
	String() string
}

type economicTreasuryInterval string

func (c economicTreasuryInterval) privateeconomicTreasuryInterval() {}
func (c economicTreasuryInterval) String() string {
	return string(c)
}

const EconomicTreasuryDaily economicTreasuryInterval = "daily"
const EconomicTreasuryWeekly economicTreasuryInterval = "weekly"
const EconomicTreasuryMonthly economicTreasuryInterval = "monthly"

type EconomicFundsInterval interface {
	privateeconomicFundsInterval()
	String() string
}

type economicFundsInterval string

func (c economicFundsInterval) privateeconomicFundsInterval() {}
func (c economicFundsInterval) String() string {
	return string(c)
}

const EconomicFundsDaily economicFundsInterval = "daily"
const EconomicFundsWeekly economicFundsInterval = "weekly"
const EconomicFundsMonthly economicFundsInterval = "monthly"

type EconomicCPIInterval interface {
	privateeconomicCPIInterval()
	String() string
}

type economicCPIInterval string

func (c economicCPIInterval) privateeconomicCPIInterval() {}
func (c economicCPIInterval) String() string {
	return string(c)
}

const EconomicCPISemiannual economicCPIInterval = "semiannual"
const EconomicCPIMonthly economicCPIInterval = "monthly"
