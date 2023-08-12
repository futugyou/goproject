package viewmodel

import "time"

type ParameterViewModel struct {
	Id        string    `bson:"_id,omitempty"`
	AccountId string    `bson:"account_id"`
	Region    string    `bson:"region"`
	Key       string    `bson:"key"`
	Value     string    `bson:"value"`
	Version   string    `bson:"version"`
	OperateAt time.Time `bson:"operate_at"`
}
