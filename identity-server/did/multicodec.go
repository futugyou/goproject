package did

import (
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

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
	s, err := NewGenericECDSAKey("ES256K")
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

var _ IVerificationMethod = (*Es256VerificationMethod)(nil)

type Es256VerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (e *Es256VerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	s, err := NewGenericECDSAKey("ES256")
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
func (e *Es256VerificationMethod) GetCrvOrSize() string {
	return "P-256"
}

// GetKeySize implements [IVerificationMethod].
func (e *Es256VerificationMethod) GetKeySize() int {
	return 33
}

// GetKty implements [IVerificationMethod].
func (e *Es256VerificationMethod) GetKty() string {
	return "EC"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (e *Es256VerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x1306"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (e *Es256VerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0x1200"
}

var _ IVerificationMethod = (*Es384VerificationMethod)(nil)

type Es384VerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (e *Es384VerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	s, err := NewGenericECDSAKey("ES384")
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
func (e *Es384VerificationMethod) GetCrvOrSize() string {
	return "P-384"
}

// GetKeySize implements [IVerificationMethod].
func (e *Es384VerificationMethod) GetKeySize() int {
	return 49
}

// GetKty implements [IVerificationMethod].
func (e *Es384VerificationMethod) GetKty() string {
	return "EC"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (e *Es384VerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x1307"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (e *Es384VerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0x1201"
}

var _ IVerificationMethod = (*JwkJcsPubVerificationMethod)(nil)

type JwkJcsPubVerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	b, err := jwk.ParseKey(publicKey)
	if err != nil {
		return nil, err
	}
	return NewJsonWebKeySecurityKey(b), nil
}

// GetCrvOrSize implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) GetCrvOrSize() string {
	return ""
}

// GetKeySize implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) GetKeySize() int {
	return 0
}

// GetKty implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) GetKty() string {
	return ""
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0xd1d603"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (j *JwkJcsPubVerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0xd1d603"
}

var _ IVerificationMethod = (*RSAVerificationMethod)(nil)

type RSAVerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (r *RSAVerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	rsa, err := NewRSASignatureKey()
	if err != nil {
		return nil, err
	}
	err = rsa.Import(publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return rsa, nil
}

// GetCrvOrSize implements [IVerificationMethod].
func (r *RSAVerificationMethod) GetCrvOrSize() string {
	return "2048+"
}

// GetKeySize implements [IVerificationMethod].
func (r *RSAVerificationMethod) GetKeySize() int {
	return 2048
}

// GetKty implements [IVerificationMethod].
func (r *RSAVerificationMethod) GetKty() string {
	return "RSA"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (r *RSAVerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x1305"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (r *RSAVerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0x1205"
}

var _ IVerificationMethod = (*X25519VerificationMethod)(nil)

type X25519VerificationMethod struct {
}

// Build implements [IVerificationMethod].
func (x *X25519VerificationMethod) Build(publicKey []byte, privateKey []byte) (IAsymmetricKey, error) {
	x25519, err := NewX25519AgreementKey()
	if err != nil {
		return nil, err
	}
	err = x25519.Import(publicKey, privateKey)
	if err != nil {
		return nil, err
	}
	return x25519, nil
}

// GetCrvOrSize implements [IVerificationMethod].
func (x *X25519VerificationMethod) GetCrvOrSize() string {
	return "X25519"
}

// GetKeySize implements [IVerificationMethod].
func (x *X25519VerificationMethod) GetKeySize() int {
	return 32
}

// GetKty implements [IVerificationMethod].
func (x *X25519VerificationMethod) GetKty() string {
	return "OKP"
}

// GetMulticodecPrivateKeyHexValue implements [IVerificationMethod].
func (x *X25519VerificationMethod) GetMulticodecPrivateKeyHexValue() string {
	return "0x1302"
}

// GetMulticodecPublicKeyHexValue implements [IVerificationMethod].
func (x *X25519VerificationMethod) GetMulticodecPublicKeyHexValue() string {
	return "0xec01"
}
