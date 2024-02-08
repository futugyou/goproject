package enums

type TimeInterval interface {
	private()
	String() string
}

type timeInterval string

func (c timeInterval) private() {}
func (c timeInterval) String() string {
	switch c {
	case T1min:
		return "1min"
	case T5min:
		return "5min"
	case T15min:
		return "15min"
	case T30min:
		return "30min"
	case T60min:
		return "60min"
	}
	return "60min"
}

const T1min timeInterval = "1min"
const T5min timeInterval = "5min"
const T15min timeInterval = "15min"
const T30min timeInterval = "30min"
const T60min timeInterval = "60min"
