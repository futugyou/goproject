package platform

type Secret struct {
	Key   string `json:"key"`   // vault aliases
	Value string `json:"value"` // vault id
}

func (s Secret) GetKey() string {
	return s.Key
}
