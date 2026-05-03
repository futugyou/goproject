package did

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/multiformats/go-multibase"
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

var _ ISignatureProof = (*Ed25519Signature2020Proof)(nil)

type Ed25519Signature2020Proof struct{}

func (e *Ed25519Signature2020Proof) GetType() string {
	return "Ed25519Signature2020"
}

func (e *Ed25519Signature2020Proof) GetVerificationMethod() string {
	return "Ed25519VerificationKey2020"
}

func (e *Ed25519Signature2020Proof) GetTransformationMethod() string {
	return "RDF"
}

func (e *Ed25519Signature2020Proof) GetHashingMethod() string {
	return "SHA256"
}

func (e *Ed25519Signature2020Proof) ComputeProof(proof *DataIntegrityProof, payload []byte, asymmetricKey IAsymmetricKey, alg string) error {
	a := jwa.NewSignatureAlgorithm(alg)
	signature, err := asymmetricKey.SignHash(payload, a)
	if err != nil {
		return fmt.Errorf("failed to sign: %w", err)
	}
	encoded, err := multibase.Encode(multibase.Base58BTC, signature)
	if err != nil {
		return err
	}

	proof.ProofValue = encoded
	return nil
}

func (e *Ed25519Signature2020Proof) GetSignature(proof *DataIntegrityProof) ([]byte, error) {
	_, decoded, err := multibase.Decode(proof.ProofValue)
	if err != nil {
		return nil, fmt.Errorf("failed to decode multibase signature: %w", err)
	}
	return decoded, nil
}
