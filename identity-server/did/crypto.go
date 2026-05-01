package did

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
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
func (k *RSASignatureKey) GetCrvOrSize() string { return "2048+" }
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

var _ IAsymmetricKey = (*X25519AgreementKey)(nil)

type X25519AgreementKey struct {
	priv *ecdh.PrivateKey
	pub  *ecdh.PublicKey
}

func NewX25519AgreementKey() (*X25519AgreementKey, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &X25519AgreementKey{priv: priv, pub: priv.PublicKey()}, nil
}

func (k *X25519AgreementKey) GetKty() string       { return "OKP" }
func (k *X25519AgreementKey) GetCrvOrSize() string { return "X25519" }
func (k *X25519AgreementKey) GetJwtAlg() string    { return "X25519" }

func (k *X25519AgreementKey) Import(publicKey []byte, privateKey []byte) error {
	if publicKey != nil {
		p, err := ecdh.X25519().NewPublicKey(publicKey)
		if err != nil {
			return err
		}
		k.pub = p
	}
	if privateKey != nil {
		p, err := ecdh.X25519().NewPrivateKey(privateKey)
		if err != nil {
			return err
		}
		k.priv = p
		k.pub = p.PublicKey()
	}
	return nil
}

func (k *X25519AgreementKey) GetPublicKey(compressed bool) []byte {
	return k.pub.Bytes()
}

func (k *X25519AgreementKey) GetPrivateKey() []byte {
	return k.priv.Bytes()
}

func (k *X25519AgreementKey) GetPublicJwk() (jwk.Key, error) {
	return jwk.Import[jwk.OKPPublicKey](k.pub)
}

func (k *X25519AgreementKey) GetPrivateJwk() (jwk.Key, error) {
	return jwk.Import[jwk.OKPPrivateKey](k.priv)
}

func (k *X25519AgreementKey) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	return nil, errors.New("X25519 does not support signing")
}

func (k *X25519AgreementKey) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	return false
}

var _ IAsymmetricKey = (*ECDSASignatureKey)(nil)

type ECDSASignatureKey struct {
	priv *ecdsa.PrivateKey
	pub  *ecdsa.PublicKey
	alg  jwa.SignatureAlgorithm
	crv  string
}

func NewGenericECDSAKey(algName string) (*ECDSASignatureKey, error) {
	var curve elliptic.Curve
	var crvName string
	alg := jwa.NewSignatureAlgorithm(algName)

	switch algName {
	case "ES256":
		curve = elliptic.P256()
		crvName = "P-256"
	case "ES256K":
		curve = secp256k1.S256()
		crvName = "secp256k1"
	case "ES384":
		curve = elliptic.P384()
		crvName = "P-384"
	case "ES512":
		curve = elliptic.P521()
		crvName = "P-521"
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algName)
	}

	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &ECDSASignatureKey{
		priv: priv,
		pub:  &priv.PublicKey,
		alg:  alg,
		crv:  crvName,
	}, nil
}

func (k *ECDSASignatureKey) GetKty() string       { return "EC" }
func (k *ECDSASignatureKey) GetCrvOrSize() string { return k.crv }
func (k *ECDSASignatureKey) GetJwtAlg() string    { return k.alg.String() }

func (k *ECDSASignatureKey) GetPublicKey(compressed bool) []byte {
	if k.pub == nil {
		return nil
	}
	ecdhPub, err := k.pub.ECDH()
	if err != nil {
		return nil
	}
	return ecdhPub.Bytes()
}

func (k *ECDSASignatureKey) GetPrivateKey() []byte {
	if k.priv == nil {
		return nil
	}
	b, err := x509.MarshalPKCS8PrivateKey(k.priv)
	if err != nil {
		return nil
	}
	return b
}

func (k *ECDSASignatureKey) GetPublicJwk() (jwk.Key, error) {
	key, err := jwk.Import[jwk.ECDSAPublicKey](k.pub)
	if err != nil {
		return nil, err
	}
	_ = key.Set(jwk.AlgorithmKey, k.alg)
	return key, nil
}

func (k *ECDSASignatureKey) GetPrivateJwk() (jwk.Key, error) {
	key, err := jwk.Import[jwk.ECDSAPrivateKey](k.priv)
	if err != nil {
		return nil, err
	}
	_ = key.Set(jwk.AlgorithmKey, k.alg)
	return key, nil
}

func (k *ECDSASignatureKey) SignHash(content []byte, alg jwa.SignatureAlgorithm) ([]byte, error) {
	if k.priv == nil {
		return nil, errors.New("missing private key")
	}

	var hash []byte
	switch alg.String() {
	case "ES256", "ES256K":
		digest := sha256.Sum256(content)
		hash = digest[:]
	case "ES384":
		digest := sha512.Sum384(content)
		hash = digest[:]
	case "ES512":
		digest := sha512.Sum512(content)
		hash = digest[:]
	default:
		return nil, fmt.Errorf("unsupported algorithm for signing: %s", alg)
	}

	r, s, err := ecdsa.Sign(rand.Reader, k.priv, hash)
	if err != nil {
		return nil, err
	}

	curveOrderBytes := (k.priv.Curve.Params().BitSize + 7) / 8
	signature := make([]byte, curveOrderBytes*2)

	r.FillBytes(signature[0:curveOrderBytes])
	s.FillBytes(signature[curveOrderBytes:])

	return signature, nil
}

func (k *ECDSASignatureKey) CheckHash(content []byte, signature []byte, alg jwa.SignatureAlgorithm) bool {
	if k.pub == nil || len(signature) == 0 {
		return false
	}

	curveOrderBytes := (k.pub.Curve.Params().BitSize + 7) / 8
	if len(signature) != curveOrderBytes*2 {
		return false
	}

	var hash []byte
	switch k.alg.String() {
	case "ES256", "ES256K":
		h := sha256.Sum256(content)
		hash = h[:]
	case "ES384":
		h := sha512.Sum384(content)
		hash = h[:]
	default:
		return false
	}

	r := new(big.Int).SetBytes(signature[:curveOrderBytes])
	s := new(big.Int).SetBytes(signature[curveOrderBytes:])

	return ecdsa.Verify(k.pub, hash, r, s)
}

func (k *ECDSASignatureKey) Import(publicKey []byte, privateKey []byte) error {
	var ecdhCurve ecdh.Curve
	switch k.alg.String() {
	case "ES256":
		ecdhCurve = ecdh.P256()
	case "ES384":
		ecdhCurve = ecdh.P384()
	case "ES512":
		ecdhCurve = ecdh.P521()
	case "ES256K":
		return k.importSecp256k1(publicKey, privateKey)
	default:
		return fmt.Errorf("unsupported alg: %s", k.alg)
	}

	if privateKey != nil {
		priv, err := ecdhCurve.NewPrivateKey(privateKey)
		if err != nil {
			return err
		}
		k.priv, err = decodeToECDSA(priv)
		if err != nil {
			return err
		}
		k.pub = &k.priv.PublicKey
	}

	if publicKey != nil && k.pub == nil {
		pub, err := ecdhCurve.NewPublicKey(publicKey)
		if err != nil {
			return err
		}
		p, err := decodeToECDSAFromPub(pub)
		if err != nil {
			return err
		}

		k.pub = p
	}

	return nil
}

func decodeToECDSA(priv *ecdh.PrivateKey) (*ecdsa.PrivateKey, error) {
	pkcs8, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}
	p, err := x509.ParsePKCS8PrivateKey(pkcs8)
	if err != nil {
		return nil, err
	}
	return p.(*ecdsa.PrivateKey), nil
}

func decodeToECDSAFromPub(pub *ecdh.PublicKey) (*ecdsa.PublicKey, error) {
	pkix, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}
	p, err := x509.ParsePKIXPublicKey(pkix)
	if err != nil {
		return nil, err
	}
	return p.(*ecdsa.PublicKey), nil
}

func (k *ECDSASignatureKey) importSecp256k1(pubBytes, privBytes []byte) error {
	curve := secp256k1.S256()

	if privBytes != nil {
		secpPriv := secp256k1.PrivKeyFromBytes(privBytes)

		secpPub := secpPriv.PubKey()
		k.priv = &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: curve,
				X:     secpPub.X(),
				Y:     secpPub.Y(),
			},
			D: new(big.Int).SetBytes(secpPriv.Serialize()),
		}
		k.pub = &k.priv.PublicKey
	}

	if pubBytes != nil && k.pub == nil {
		secpPub, err := secp256k1.ParsePubKey(pubBytes)
		if err != nil {
			return err
		}
		k.pub = &ecdsa.PublicKey{
			Curve: curve,
			X:     secpPub.X(),
			Y:     secpPub.Y(),
		}
	}

	return nil
}
