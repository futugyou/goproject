package entity

type KeyValueEntity struct {
	Id    string `bson:"_id,omitempty"`
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (KeyValueEntity) GetType() string {
	return "keyvalues"
}
