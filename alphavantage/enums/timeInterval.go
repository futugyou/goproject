package enums

type TimeInterval interface {
	privateTimeInterval()
	String() string
}

type timeInterval string

func (c timeInterval) privateTimeInterval() {}
func (c timeInterval) String() string {
	return string(c)
}

const T1min timeInterval = "1min"
const T5min timeInterval = "5min"
const T15min timeInterval = "15min"
const T30min timeInterval = "30min"
const T60min timeInterval = "60min"
