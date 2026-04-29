package did

type IMulticodecSerializer interface {
	Deserialize(publicKey string, privateKey string) (IAsymmetricKey, error)
	SerializePublicKey(signatureKey IAsymmetricKey) string
	SerializePrivateKey(signatureKey IAsymmetricKey) string
}

type IVerificationMethod interface {
	GetMulticodecPublicKeyHexValue() string
	GetMulticodecPrivateKeyHexValue() string
	GetKeySize() int
	GetKty() string
	GetCrvOrSize() string
	Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error)
}
