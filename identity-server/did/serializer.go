package did

import (
	"errors"
	"strings"

	"github.com/multiformats/go-multibase"
)

func Encode(payload []byte) (string, error) {
	return multibase.Encode(multibase.Base58BTC, payload)
}

func Decode(input string) ([]byte, error) {
	_, decoded, err := multibase.Decode(input)
	return decoded, err
}

type IMulticodecSerializer interface {
	Deserialize(publicKey string, privateKey string) (IAsymmetricKey, error)
	SerializePublicKey(signatureKey IAsymmetricKey) (string, error)
	SerializePrivateKey(signatureKey IAsymmetricKey) (string, error)
}

var _ IMulticodecSerializer = (*MulticodecSerializer)(nil)

type MulticodecSerializer struct {
	verificationMethods []IVerificationMethod
}

func NewMulticodecSerializer(methods []IVerificationMethod) *MulticodecSerializer {
	return &MulticodecSerializer{verificationMethods: methods}
}

// Deserialize string keys
func (m *MulticodecSerializer) Deserialize(publicKey string, privateKey string) (IAsymmetricKey, error) {
	publicKeyPayload, err := Decode(publicKey)
	if err != nil {
		return nil, err
	}
	var privateKeyPayload []byte
	if privateKey != "" {
		privateKeyPayload, err = Decode(privateKey)
		if err != nil {
			return nil, err
		}
	}
	return m.DeserializeBytes(publicKeyPayload, privateKeyPayload)
}

// Deserialize byte slices
func (m *MulticodecSerializer) DeserializeBytes(publicKeyPayload []byte, privateKeyPayload []byte) (IAsymmetricKey, error) {
	if publicKeyPayload == nil {
		return nil, errors.New("publicKeyPayload cannot be nil")
	}

	hex := ToHex(publicKeyPayload, true)
	var verificationMethod IVerificationMethod
	for _, v := range m.verificationMethods {
		if strings.HasPrefix(hex, v.GetMulticodecPublicKeyHexValue()) {
			verificationMethod = v
			break
		}
	}
	if verificationMethod == nil {
		return nil, errors.New("public key; either the multicodec is invalid or the verification method is not supported")
	}

	var privateKey []byte
	if privateKeyPayload != nil {
		hexPrivateKey := ToHex(privateKeyPayload, true)
		if !strings.HasPrefix(hexPrivateKey, verificationMethod.GetMulticodecPrivateKeyHexValue()) {
			return nil, errors.New("private key; either the multicodec is invalid or the verification method is not supported")
		}
		privateKeyHeader, _ := HexToByteArray(verificationMethod.GetMulticodecPrivateKeyHexValue())
		privateKey = privateKeyPayload[len(privateKeyHeader):]
	}

	publicKeyHeader, _ := HexToByteArray(verificationMethod.GetMulticodecPublicKeyHexValue())
	publicKey := publicKeyPayload[len(publicKeyHeader):]

	return verificationMethod.Build(publicKey, privateKey)
}

// SerializePublicKey
func (m *MulticodecSerializer) SerializePublicKey(signatureKey IAsymmetricKey) (string, error) {
	var verificationMethod IVerificationMethod
	for _, v := range m.verificationMethods {
		if v.GetKty() == signatureKey.GetKty() && v.GetCrvOrSize() == signatureKey.GetCrvOrSize() {
			verificationMethod = v
			break
		}
	}

	if verificationMethod == nil {
		return "", errors.New("verification method not found")
	}
	d, _ := HexToByteArray(verificationMethod.GetMulticodecPublicKeyHexValue())
	publicKey := append(d, signatureKey.GetPublicKey(false)...)
	return Encode(publicKey)
}

// SerializePrivateKey
func (m *MulticodecSerializer) SerializePrivateKey(signatureKey IAsymmetricKey) (string, error) {
	var verificationMethod IVerificationMethod
	for _, v := range m.verificationMethods {
		if v.GetKty() == signatureKey.GetKty() && v.GetCrvOrSize() == signatureKey.GetCrvOrSize() {
			verificationMethod = v
			break
		}
	}

	if verificationMethod == nil {
		return "", errors.New("verification method not found")
	}

	d, _ := HexToByteArray(verificationMethod.GetMulticodecPrivateKeyHexValue())
	privateKey := append(d, signatureKey.GetPrivateKey()...)
	return Encode(privateKey)
}
