package enums

type Maturity interface {
	private()
	String() string
}

type maturity string

func (c maturity) private() {}
func (c maturity) String() string {
	return string(c)
}

const M3month maturity = "3month"
const M2year maturity = "2year"
const M5year maturity = "5year"
const M7year maturity = "7year"
const M10year maturity = "10year"
const M30year maturity = "30year "
