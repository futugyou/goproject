package did

import (
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
