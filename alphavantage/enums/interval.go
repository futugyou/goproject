package enums

type Interval interface {
	privateInterval()
	String() string
}

type interval string

func (c interval) privateInterval() {}
func (c interval) String() string {
	return string(c)
}

const I1min interval = "1min"
const I5min interval = "5min"
const I15min interval = "15min"
const I30min interval = "30min"
const I60min interval = "60min"
const IDaily interval = "DAILY"
const IWeekly interval = "WEEKLY"
const IMonthly interval = "MONTHLY"
