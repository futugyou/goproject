package stockSeries

import "time"

type StockSeriesEntity struct {
	Id     string    `bson:"_id,omitempty"`
	Symbol string    `bson:"symbol"`
	Time   time.Time `bson:"time"`
	Open   float64   `bson:"open"`
	High   float64   `bson:"high"`
	Low    float64   `bson:"low"`
	Close  float64   `bson:"close"`
	Volume float64   `bson:"volume"`
}

func (StockSeriesEntity) GetType() string {
	return "stock-series"
}

type StockSeriesConfigEntity struct {
	Id     string `bson:"_id,omitempty"`
	Month  string `bson:"month"`
	Filter string `bson:"filter"`
}

func (StockSeriesConfigEntity) GetType() string {
	return "stock-series-config"
}
