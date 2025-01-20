package platform

type Property struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

func (p Property) GetKey() string {
	return p.Key
}
