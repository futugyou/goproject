package platform

type PropertyInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (p PropertyInfo) GetKey() string {
	return p.Key
}
