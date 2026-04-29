package did

import (
	"crypto"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/v4/jwk"
)

func ToHex(value []byte, prefix bool) string {
	res := hex.EncodeToString(value)
	if prefix {
		return "0x" + res
	}
	return res
}

func HexToByteArray(value string) ([]byte, error) {
	value = strings.TrimPrefix(value, "0x")

	if len(value)%2 != 0 {
		value = "0" + value
	}

	return hex.DecodeString(value)
}

func ComputeJWKThumbprint(key jwk.Key) ([]byte, error) {
	return key.Thumbprint(crypto.SHA256)
}

type OKPThumbprint struct {
	Crv string `json:"crv"`
	Kty string `json:"kty"`
	X   string `json:"x"`
}

func ComputeOKPThumbprint(key jwk.Key) ([]byte, error) {
	var crv, kty, x any
	var ok bool

	if crv, ok = key.Field("crv"); !ok {
		return nil, fmt.Errorf("crv can not be null")
	}
	if kty, ok = key.Field("kty"); !ok {
		return nil, fmt.Errorf("kty can not be null")
	}
	if x, ok = key.Field("x"); !ok {
		return nil, fmt.Errorf("x can not be null")
	}

	thumb := OKPThumbprint{
		Crv: crv.(string),
		Kty: kty.(string),
		X:   x.(string),
	}

	canonicalJwk, err := json.Marshal(thumb)
	if err != nil {
		return nil, err
	}

	h := sha256.Sum256(canonicalJwk)
	return h[:], nil
}
