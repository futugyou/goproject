package expected

import "time"

type ExpectedEntity struct {
	Id               string    `bson:"_id"`
	Symbol           string    `bson:"symbol"`
	Name             string    `bson:"name"`
	ReportDate       time.Time `bson:"reportDate"`
	FiscalDateEnding time.Time `bson:"fiscalDateEnding"`
	Estimate         float64   `bson:"estimate"`
	Currency         string    `bson:"currency"`
}

func (ExpectedEntity) GetType() string {
	return "Expecteds"
}
