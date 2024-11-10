package platform

type Secret struct {
	Key            string `json:"key" bson:"key"`                           // vault aliases
	Value          string `json:"value" bson:"value"`                       // vault id
	VaultKey       string `json:"vault_key" bson:"vault_key"`               // vault key
	VaultMaskValue string `json:"vault_mask_value" bson:"vault_mask_value"` // vault MaskValue
}

func (s Secret) GetKey() string {
	return s.Key
}
