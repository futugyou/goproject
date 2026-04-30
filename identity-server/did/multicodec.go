package did

import (
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

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

var _ IVerificationMethod = (*Ec256VerificationMethod)(nil)

type Ec256VerificationMethod struct {
}

// CheckHash implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	panic("unimplemented")
}

// GetJwtAlg implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) GetJwtAlg() string {
	panic("unimplemented")
}

// GetPrivateJwk implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) GetPrivateJwk() (jwk.Key, error) {
	panic("unimplemented")
}

// GetPrivateKey implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) GetPrivateKey() []byte {
	panic("unimplemented")
}

// GetPublicJwk implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) GetPublicJwk() (jwk.Key, error) {
	panic("unimplemented")
}

// GetPublicKey implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) GetPublicKey(compressed bool) []byte {
	panic("unimplemented")
}

// Import implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) Import(publicKey []byte, privateKey []byte) error {
	panic("unimplemented")
}

// SignHash implements [IAsymmetricKey].
func (e *Ec256VerificationMethod) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	panic("unimplemented")
}

// Build implements [IVerificationMethod].
func (e *Ec256VerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	panic("unimplemented")
}

// GetCrvOrSize implements [IVerificationMethod].
func (e *Ec256VerificationMethod) GetCrvOrSize() string {
	panic("unimplemented")
}

// GetKeySize implements [IVerificationMethod].
func (e *Ec256VerificationMethod) GetKeySize() int {
	return 0
}

// GetKty implements [IVerificationMethod].
func (e *Ec256VerificationMethod) GetKty() string {
	panic("unimplemented")
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (e *Ec256VerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return ""
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (e *Ec256VerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return ""
}
