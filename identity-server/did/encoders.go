package did

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

type SignatureKeyEncodingTypes uint32

const (
	BASE58    SignatureKeyEncodingTypes = 1 << iota // 1 (2^0)
	HEX                                             // 2 (2^1)
	MULTIBASE                                       // 4 (2^2)
	JWK                                             // 8 (2^3)
)

type IVerificationMethodEncoding interface {
	Encode(encodeType string, did string, controller string, key IAsymmetricKey, encoding SignatureKeyEncodingTypes, includePrivateKey bool) (*DidDocumentVerificationMethod, error)
	Decode(verificationMethod DidDocumentVerificationMethod) (IAsymmetricKey, error)
	GetStandards() []IVerificationMethodStandard
}

type IVerificationMethodStandard interface {
	GetJSONLDContext() string
	GetType() string
	GetDefaultEncoding() SignatureKeyEncodingTypes
	GetSupportedEncoding() SignatureKeyEncodingTypes
	GetSupportedCurves() []string
	BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string
}

type VerificationMethodEncoding struct {
	standards            []IVerificationMethodStandard
	multicodecSerializer IMulticodecSerializer
	verificationMethods  []IVerificationMethod
}

func NewVerificationMethodEncoding(standards []IVerificationMethodStandard, multicodecSerializer IMulticodecSerializer, verificationMethods []IVerificationMethod) *VerificationMethodEncoding {
	return &VerificationMethodEncoding{
		standards:            standards,
		multicodecSerializer: multicodecSerializer,
		verificationMethods:  verificationMethods,
	}
}

func (v *VerificationMethodEncoding) Encode(encodeType string, did string, controller string, key IAsymmetricKey, encoding SignatureKeyEncodingTypes, includePrivateKey bool) (*DidDocumentVerificationMethod, error) {
	var privateKey []byte
	if includePrivateKey {
		privateKey = key.GetPrivateKey()
		if len(privateKey) == 0 {
			return nil, fmt.Errorf("private key is required but not provided")
		}
	}
	index := slices.IndexFunc(v.standards, func(stn IVerificationMethodStandard) bool {
		return stn.GetType() == encodeType
	})

	standard := v.standards[index]
	result := &DidDocumentVerificationMethod{
		Type:       standard.GetType(),
		Controller: controller,
	}

	switch encoding {
	case BASE58:
		result.PublicKeyBase58 = base58.Encode(key.GetPublicKey(false))
		if includePrivateKey {
			result.PrivateKeyBase58 = base58.Encode(privateKey)
		}
	case HEX:
		result.PublicKeyHex = ToHex(key.GetPublicKey(true), false)
		if includePrivateKey {
			result.PrivateKeyHex = ToHex(key.GetPrivateKey(), false)
		}
	case MULTIBASE:
		result.PublicKeyMultibase = v.multicodecSerializer.SerializePublicKey(key)
		if includePrivateKey {
			result.SecretKeyMultibase = v.multicodecSerializer.SerializePrivateKey(key)
		}
	case JWK:
		publicKey, err := key.GetPublicJwk()
		if err != nil {
			return nil, err
		}

		keyBytes, err := json.Marshal(publicKey)
		if err != nil {
			return nil, err
		}

		result.PublicKeyJwk = keyBytes
		if includePrivateKey {
			privateKey, err := key.GetPrivateJwk()
			if err != nil {
				return nil, err
			}
			keyBytes, err := json.Marshal(privateKey)
			if err != nil {
				return nil, err
			}
			result.PrivateKeyJwk = keyBytes
		}
	}

	return result, nil
}

func (v *VerificationMethodEncoding) Decode(verificationMethod DidDocumentVerificationMethod) (IAsymmetricKey, error) {
	index := slices.IndexFunc(v.standards, func(stn IVerificationMethodStandard) bool {
		return stn.GetType() == verificationMethod.Type
	})

	standard := v.standards[index]
	var publicKey []byte
	switch standard.GetDefaultEncoding() {
	case BASE58:
		publicKey = base58.Decode(verificationMethod.PublicKeyBase58)
	case HEX:
		publicKey, _ = HexToByteArray(verificationMethod.PublicKeyHex)
	case MULTIBASE:
		return v.multicodecSerializer.Deserialize(verificationMethod.PublicKeyMultibase, verificationMethod.SecretKeyMultibase)
	case JWK:
		key, err := jwk.ParseKey(verificationMethod.PublicKeyJwk)
		if err != nil {
			return nil, err
		}
		crv, bool := key.Field(jwk.ECDSACrvKey)
		if !bool {
			return nil, fmt.Errorf("invalid JWK format")
		}
		for _, vm := range v.verificationMethods {
			if vm.GetKty() == key.KeyType().String() && vm.GetCrvOrSize() == crv {
				return vm.Build(verificationMethod.PublicKeyJwk, verificationMethod.PrivateKeyJwk)
			}
		}
	}

	for _, curve := range standard.GetSupportedCurves() {
		for _, me := range v.verificationMethods {
			if me.GetCrvOrSize() == curve {
				return me.Build(publicKey, nil)
			}
		}
	}

	return nil, fmt.Errorf("unsupported encoding type or curve")
}

func (v *VerificationMethodEncoding) GetStandards() []IVerificationMethodStandard {
	return v.standards
}
