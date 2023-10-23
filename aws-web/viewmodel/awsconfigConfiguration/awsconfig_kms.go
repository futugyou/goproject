package awsconfigConfiguration

type KMSConfiguration struct {
	KeyID                 string        `json:"keyId"`
	Arn                   string        `json:"arn"`
	CreationDate          int64         `json:"creationDate"`
	Enabled               bool          `json:"enabled"`
	Description           string        `json:"description"`
	KeyUsage              string        `json:"keyUsage"`
	KeyState              string        `json:"keyState"`
	Origin                string        `json:"origin"`
	KeyManager            string        `json:"keyManager"`
	CustomerMasterKeySpec string        `json:"customerMasterKeySpec"`
	KeySpec               string        `json:"keySpec"`
	EncryptionAlgorithms  []string      `json:"encryptionAlgorithms"`
	SigningAlgorithms     []interface{} `json:"signingAlgorithms"`
	MultiRegion           bool          `json:"multiRegion"`
	MACAlgorithms         []interface{} `json:"macAlgorithms"`
	AwsaccountID          string        `json:"awsaccountId"`
}
