package platform

type Secret struct {
	Key            string `bson:"key"`              // vault aliases
	Value          string `bson:"value"`            // vault id
	VaultKey       string `bson:"vault_key"`        // vault key
	VaultMaskValue string `bson:"vault_mask_value"` // vault MaskValue
}

func (s Secret) GetKey() string {
	return s.Key
}
