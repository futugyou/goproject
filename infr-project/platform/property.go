package platform

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p Property) GetKey() string {
	return p.Key
}
