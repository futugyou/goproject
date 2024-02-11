package enums

type CommoditiesInterval interface {
	privateCommoditiesInterval()
	String() string
}

type commoditiesInterval string

func (c commoditiesInterval) privateCommoditiesInterval() {}
func (c commoditiesInterval) String() string {
	return string(c)
}

const CommoditiesDaily commoditiesInterval = "DAILY"
const CommoditiesWeekly commoditiesInterval = "WEEKLY"
const CommoditiesMonthly commoditiesInterval = "MONTHLY"

type CommoditiesInterval2 interface {
	privateCommoditiesInterval2()
	String() string
}

type commoditiesInterval2 string

func (c commoditiesInterval2) privateCommoditiesInterval2() {}
func (c commoditiesInterval2) String() string {
	return string(c)
}

const CommoditiesMonthly2 commoditiesInterval2 = "monthly"
const CommoditiesQuarterly2 commoditiesInterval2 = "quarterly"
const CommoditiesAnnual2 commoditiesInterval2 = "annual "
