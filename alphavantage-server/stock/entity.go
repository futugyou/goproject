package stock

type StockEntity struct {
	Id          string  `bson:"_id,omitempty"`
	Symbol      string  `bson:"symbol"`
	Name        string  `bson:"name"`
	Type        string  `bson:"type"`
	Region      string  `bson:"region"`
	MarketOpen  string  `bson:"marketOpen"`
	MarketClose string  `bson:"marketClose"`
	Timezone    string  `bson:"timezone"`
	Currency    string  `bson:"currency"`
	MatchScore  float64 `bson:"matchScore"`
}

func (StockEntity) GetType() string {
	return "stocks"
}
