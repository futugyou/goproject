package entity

type ParameterSearchFilter struct {
	AccountId string `bson:"account_id"`
	Region    string `bson:"region"`
	Key       string `bson:"key"`
}

type ParameterEntity struct {
	Id        string `bson:"_id,omitempty"`
	AccountId string `bson:"account_id"`
	Region    string `bson:"region"`
	Key       string `bson:"key"`
	Value     string `bson:"value"`
	Version   string `bson:"version"`
	OperateAt int64  `bson:"operate_at"`
}

func (ParameterEntity) GetType() string {
	return "parameters"
}

type ParameterLogEntity struct {
	Id        string `bson:"_id,omitempty"`
	AccountId string `bson:"account_id"`
	Region    string `bson:"region"`
	Key       string `bson:"key"`
	Value     string `bson:"value"`
	Version   string `bson:"version"`
	OperateAt int64  `bson:"operate_at"`
}

func (ParameterLogEntity) GetType() string {
	return "parameter_logs"
}
