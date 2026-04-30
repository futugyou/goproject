package did

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

type IAsymmetricKey interface {
	GetKty() string
	GetCrvOrSize() string
	GetJwtAlg() string
	Import(publicKey []byte, privateKey []byte) error
	GetPublicKey(compressed bool) []byte
	GetPrivateKey() []byte
	GetPublicJwk() (jwk.Key, error)
	GetPrivateJwk() (jwk.Key, error)
	SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error)
	CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool
}

var _ IAsymmetricKey = (*Ec256VerificationMethod)(nil)

type Ed25519SignatureKey struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func NewEd25519SignatureKey() *Ed25519SignatureKey {
	return &Ed25519SignatureKey{}
}

func Ed25519SignatureKeyGenerate() (*Ed25519SignatureKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Ed25519SignatureKey{
		publicKey:  pub,
		privateKey: priv,
	}, nil
}

func Ed25519SignatureKeyFrom(publicKey []byte, privateKey []byte) (*Ed25519SignatureKey, error) {
	key := &Ed25519SignatureKey{}
	err := key.Import(publicKey, privateKey)
	return key, err
}

func (k *Ed25519SignatureKey) GetKty() string {
	return "OKP"
}

func (k *Ed25519SignatureKey) GetCrvOrSize() string {
	return "Ed25519"
}

func (k *Ed25519SignatureKey) GetJwtAlg() string {
	return "EdDSA"
}

func (k *Ed25519SignatureKey) Import(publicKey []byte, privateKey []byte) error {
	if publicKey != nil {
		if len(publicKey) != ed25519.PublicKeySize {
			return fmt.Errorf("public key must have %d bytes", ed25519.PublicKeySize)
		}
		k.publicKey = make([]byte, ed25519.PublicKeySize)
		copy(k.publicKey, publicKey)
	}

	if privateKey != nil {
		if len(privateKey) == ed25519.SeedSize {
			k.privateKey = ed25519.NewKeyFromSeed(privateKey)
			k.publicKey = k.privateKey.Public().(ed25519.PublicKey)
		} else if len(privateKey) == ed25519.PrivateKeySize {
			k.privateKey = make([]byte, ed25519.PrivateKeySize)
			copy(k.privateKey, privateKey)
			k.publicKey = k.privateKey.Public().(ed25519.PublicKey)
		} else {
			return fmt.Errorf("private key must have %d or %d bytes", ed25519.SeedSize, ed25519.PrivateKeySize)
		}
	}
	return nil
}

func (k *Ed25519SignatureKey) GetPublicKey(compressed bool) []byte {
	return k.publicKey
}

func (k *Ed25519SignatureKey) GetPrivateKey() []byte {
	return k.privateKey
}

func (k *Ed25519SignatureKey) GetPublicJwk() (jwk.Key, error) {
	if k.publicKey == nil {
		return nil, errors.New("there is no public key")
	}
	key, err := jwk.Import[jwk.OKPPublicKey](k.publicKey)
	if err != nil {
		return nil, err
	}
	_ = key.Set(jwk.AlgorithmKey, jwa.EdDSA)
	_ = key.Set(jwk.KeyUsageKey, jwk.ForSignature)
	return key, nil
}

func (k *Ed25519SignatureKey) GetPrivateJwk() (jwk.Key, error) {
	if k.privateKey == nil {
		return nil, errors.New("there is no private key")
	}
	key, err := jwk.Import[jwk.OKPPrivateKey](k.privateKey)
	if err != nil {
		return nil, err
	}
	_ = key.Set(jwk.AlgorithmKey, jwa.EdDSA)
	_ = key.Set(jwk.KeyUsageKey, jwk.ForSignature)
	return key, nil
}

func (k *Ed25519SignatureKey) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	if k.privateKey == nil {
		return nil, errors.New("there is no private key")
	}
	return ed25519.Sign(k.privateKey, content), nil
}

func (k *Ed25519SignatureKey) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	if k.publicKey == nil {
		return false
	}
	return ed25519.Verify(k.publicKey, content, signature)
}

var _ IAsymmetricKey = (*JsonWebKeySecurityKey)(nil)

type JsonWebKeySecurityKey struct {
	innerKey jwk.Key
}

func NewJsonWebKeySecurityKey(key jwk.Key) *JsonWebKeySecurityKey {
	return &JsonWebKeySecurityKey{
		innerKey: key,
	}
}

func (k *JsonWebKeySecurityKey) GetKty() string {
	return k.innerKey.KeyType().String()
}

func (k *JsonWebKeySecurityKey) GetCrvOrSize() string {
	var crv string
	crv, err := jwk.Get[string](k.innerKey, "crv")
	if err == nil {
		return crv
	}
	return ""
}

func (k *JsonWebKeySecurityKey) GetJwtAlg() string {
	alg, ok := k.innerKey.Algorithm()
	if ok {
		return alg.String()
	}
	return ""
}

func (k *JsonWebKeySecurityKey) GetPublicJwk() (jwk.Key, error) {
	return jwk.PublicKeyOf(k.innerKey)
}

func (k *JsonWebKeySecurityKey) GetPrivateJwk() (jwk.Key, error) {
	return k.innerKey, nil
}

func (k *JsonWebKeySecurityKey) GetPublicKey(compressed bool) []byte {
	pub, _ := jwk.PublicKeyOf(k.innerKey)
	buf, _ := json.Marshal(pub)
	return buf
}

func (k *JsonWebKeySecurityKey) GetPrivateKey() []byte {
	buf, _ := json.Marshal(k.innerKey)
	return buf
}

func (k *JsonWebKeySecurityKey) Import(publicKey []byte, privateKey []byte) error {
	return errors.New("import not supported in JsonWebKeySecurityKey")
}

func (k *JsonWebKeySecurityKey) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	return nil, errors.New("not implemented")
}

func (k *JsonWebKeySecurityKey) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	return false
}

var _ IAsymmetricKey = (*RSASignatureKey)(nil)

type RSASignatureKey struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
}

func NewRSASignatureKey() (*RSASignatureKey, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	return &RSASignatureKey{priv: priv, pub: &priv.PublicKey}, nil
}

func (k *RSASignatureKey) GetKty() string       { return "RSA" }
func (k *RSASignatureKey) GetCrvOrSize() string { return "2048" }
func (k *RSASignatureKey) GetJwtAlg() string    { return "RS256" }

func (k *RSASignatureKey) Import(publicKey []byte, privateKey []byte) error {
	if publicKey != nil {
		p, err := x509.ParsePKCS1PublicKey(publicKey)
		if err != nil {
			return err
		}
		k.pub = p
	}
	if privateKey != nil {
		p, err := x509.ParsePKCS1PrivateKey(privateKey)
		if err != nil {
			return err
		}
		k.priv = p
		k.pub = &p.PublicKey
	}
	return nil
}

func (k *RSASignatureKey) GetPublicKey(compressed bool) []byte {
	return x509.MarshalPKCS1PublicKey(k.pub)
}

func (k *RSASignatureKey) GetPrivateKey() []byte {
	if k.priv == nil {
		return nil
	}
	return x509.MarshalPKCS1PrivateKey(k.priv)
}

func (k *RSASignatureKey) GetPublicJwk() (jwk.Key, error) {
	return jwk.Import[jwk.RSAPublicKey](k.pub)
}

func (k *RSASignatureKey) GetPrivateJwk() (jwk.Key, error) {
	return jwk.Import[jwk.RSAPrivateKey](k.priv)
}

func (k *RSASignatureKey) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	return rsa.SignPKCS1v15(rand.Reader, k.priv, crypto.SHA256, content)
}

func (k *RSASignatureKey) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	err := rsa.VerifyPKCS1v15(k.pub, crypto.SHA256, content, signature)
	return err == nil
}
