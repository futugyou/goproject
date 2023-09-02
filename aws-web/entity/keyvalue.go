package entity

type KeyValueEntity struct {
	Id    string `bson:"_id,omitempty"`
	Value string `bson:"value"`
}

func (KeyValueEntity) GetType() string {
	return "keyvalues"
}
