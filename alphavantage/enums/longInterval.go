package enums

type LongInterval interface {
	private()
	String() string
}

type longInterval string

func (c longInterval) private() {}
func (c longInterval) String() string {
	return string(c)
}

const LDaily longInterval = "DAILY"
const LWeekly longInterval = "WEEKLY"
const LMonthly longInterval = "MONTHLY"
