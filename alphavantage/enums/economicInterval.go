package enums

type EconomicInterval interface {
	privateeconomicInterval()
	String() string
}

type economicInterval string

func (c economicInterval) privateeconomicInterval() {}
func (c economicInterval) String() string {
	return string(c)
}

const EconomicQuarterly economicInterval = "quarterly"
const EconomicAnnual economicInterval = "annual"
