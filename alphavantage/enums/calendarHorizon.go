package enums

type CalendarHorizon interface {
	privateCalendarHorizon()
	String() string
}

type calendarHorizon string

func (c calendarHorizon) privateCalendarHorizon() {}
func (c calendarHorizon) String() string {
	return string(c)
}

const H3month calendarHorizon = "3month"
const H6month calendarHorizon = "6month"
const H12month calendarHorizon = "12month"
