package viewmodel

import "time"

type ParameterViewModel struct {
	Id        string    `bson:"_id,omitempty"`
	AccountId string    `bson:"account_id"`
	Region    string    `bson:"region"`
	Key       string    `bson:"key"`
	Value     string    `bson:"value,omitempty"`
	Version   string    `bson:"version"`
	NeedSync  bool      `bson:"need_sync,omitempty"`
	OperateAt time.Time `bson:"operate_at"`
}
