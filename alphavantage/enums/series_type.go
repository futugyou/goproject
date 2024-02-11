package enums

type SeriesType interface {
	privateSeriesType()
	String() string
}

type seriesType string

func (c seriesType) privateSeriesType() {}
func (c seriesType) String() string {
	return string(c)
}

const TechnicalSeriesTypeClose seriesType = "close"
const TechnicalSeriesTypeOpen seriesType = "open"
const TechnicalSeriesTypeHigh seriesType = "high"
const TechnicalSeriesTypeLow seriesType = "low"
