package enums

type OHLC interface {
	private()
	String() string
}

type ohlc string

func (c ohlc) private() {}
func (c ohlc) String() string {
	return string(c)
}

const Close ohlc = "close"
const Open ohlc = "open"
const High ohlc = "high"
const Low ohlc = "low"
