package base

import "time"

type BaseDataEntity struct {
	Id       string    `bson:"_id"`
	Symbol   string    `bson:"symbol"`
	RunCount float64   `bson:"run-count"`
	RunDate  time.Time `bson:"run-date"`
}

func (BaseDataEntity) GetType() string {
	return "base-datas"
}
