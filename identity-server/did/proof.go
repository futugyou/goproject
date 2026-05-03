package did

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/v4/jwa"
)

type ISignatureProof interface {
	GetType() string
	GetVerificationMethod() string
	GetTransformationMethod() string
	GetHashingMethod() string
	ComputeProof(proof *DataIntegrityProof, payload []byte, asymmetricKey IAsymmetricKey, alg string) error
	GetSignature(proof *DataIntegrityProof) ([]byte, error)
}

var _ ISignatureProof = (*JsonWebSignature2020Proof)(nil)

type JsonWebSignature2020Proof struct{}

// ComputeProof implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) ComputeProof(proof *DataIntegrityProof, payload []byte, asymmetricKey IAsymmetricKey, alg string) error {
	a := jwa.NewSignatureAlgorithm(alg)
	signature, err := asymmetricKey.SignHash(payload, a)
	if err != nil {
		return fmt.Errorf("failed to sign hash: %w", err)
	}

	header, err := j.buildJwtHeader()
	if err != nil {
		return err
	}

	jwtSignature := base64.RawURLEncoding.EncodeToString(signature)

	proof.Jws = fmt.Sprintf("%s..%s", header, jwtSignature)
	return nil
}

// GetHashingMethod implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) GetHashingMethod() string {
	return "SHA256"
}

// GetSignature implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) GetSignature(proof *DataIntegrityProof) ([]byte, error) {
	parts := strings.Split(proof.Jws, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWS format in proof")
	}

	jwtSignature := parts[2]

	return base64.RawURLEncoding.DecodeString(jwtSignature)
}

// GetTransformationMethod implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) GetTransformationMethod() string {
	return "RDF"
}

// GetType implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) GetType() string {
	return "JsonWebSignature2020"
}

// GetVerificationMethod implements [ISignatureProof].
func (j *JsonWebSignature2020Proof) GetVerificationMethod() string {
	return "JsonWebKey2020"
}

func (j *JsonWebSignature2020Proof) buildJwtHeader() (string, error) {
	headerObj := map[string]interface{}{
		"alg":  "EdDSA",
		"b64":  false,
		"crit": []string{"b64"},
	}

	jsonBytes, err := json.Marshal(headerObj)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(jsonBytes), nil
}
