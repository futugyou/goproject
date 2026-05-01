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

var _ IVerificationMethod = (*Ed25519VerificationMethod)(nil)

type Ed25519VerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	return Ed25519SignatureKeyFrom(publicKey, privateKey)
}

// GetCrvOrSize implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) GetCrvOrSize() string {
	return "Ed25519"
}

// GetKeySize implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) GetKeySize() int {
	return 32
}

// GetKty implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) GetKty() string {
	return "OKP"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x8026"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (e *Ed25519VerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0xed01"
}

var _ IVerificationMethod = (*Es256KVerificationMethod)(nil)

type Es256KVerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (e *Es256KVerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	s, err := NewGenericECDSAKey("EC")
	if err != nil {
		return nil, err
	}
	err = s.Import(publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return s, err
}

// GetCrvOrSize implements [IVerificationMethod].
func (e *Es256KVerificationMethod) GetCrvOrSize() string {
	return "secp256k1"
}

// GetKeySize implements [IVerificationMethod].
func (e *Es256KVerificationMethod) GetKeySize() int {
	return 33
}

// GetKty implements [IVerificationMethod].
func (e *Es256KVerificationMethod) GetKty() string {
	return "EC"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (e *Es256KVerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x1301"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (e *Es256KVerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0xe7"
}
