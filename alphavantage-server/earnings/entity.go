package earnings

type EarningsEntity struct {
	Id                 string `bson:"_id"`
	Symbol             string `bson:"symbol"`
	DataType           string `bson:"dataType"`
	FiscalDateEnding   string `bson:"fiscalDateEnding"`
	ReportedDate       string `bson:"reportedDate"`
	ReportedEPS        string `bson:"reportedEPS"`
	EstimatedEPS       string `bson:"estimatedEPS"`
	Surprise           string `bson:"surprise"`
	SurprisePercentage string `bson:"surprisePercentage"`
}

func (EarningsEntity) GetType() string {
	return "earningss"
}
