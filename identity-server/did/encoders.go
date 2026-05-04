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
		result.PublicKeyMultibase, _ = v.multicodecSerializer.SerializePublicKey(key)
		if includePrivateKey {
			result.SecretKeyMultibase, _ = v.multicodecSerializer.SerializePrivateKey(key)
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

func GetAllVerificationMethodStandards() []IVerificationMethodStandard {
	methods := []IVerificationMethod{
		&Ed25519VerificationMethod{},
		&Es256KVerificationMethod{},
		&Es256VerificationMethod{},
		&Es384VerificationMethod{},
		&X25519VerificationMethod{},
		&RSAVerificationMethod{},
		&JwkJcsPubVerificationMethod{},
	}
	multicodecSerializer := NewMulticodecSerializer(methods)

	return []IVerificationMethodStandard{
		&EcdsaSecp256k1RecoveryMethod2020Standard{},
		&EcdsaSecp256k1VerificationKey2019Standard{},
		NewEd25519VerificationKey2018Standard(multicodecSerializer),
		&Ed25519VerificationKey2020Standard{},
		&JsonWebKey2020Standard{},
		NewX25519KeyAgreementKey2019Standard(multicodecSerializer),
		&X25519KeyAgreementKey2020Standard{},
	}
}

var _ IVerificationMethodStandard = (*Ed25519VerificationKey2020Standard)(nil)

type Ed25519VerificationKey2020Standard struct {
}

// BuildId implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	return verificationMethod.PublicKeyMultibase
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return MULTIBASE
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/ed25519-2020/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) GetSupportedCurves() []string {
	return []string{"Ed25519"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return MULTIBASE
}

// GetType implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2020Standard) GetType() string {
	return "Ed25519VerificationKey2020"
}

var _ IVerificationMethodStandard = (*EcdsaSecp256k1RecoveryMethod2020Standard)(nil)

type EcdsaSecp256k1RecoveryMethod2020Standard struct {
}

// BuildId implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	return ""
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return JWK
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/secp256k1recovery-2020/v2"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) GetSupportedCurves() []string {
	return []string{"secp256k1"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return HEX | JWK
}

// GetType implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1RecoveryMethod2020Standard) GetType() string {
	return "EcdsaSecp256k1RecoveryMethod2020"
}

var _ IVerificationMethodStandard = (*EcdsaSecp256k1VerificationKey2019Standard)(nil)

type EcdsaSecp256k1VerificationKey2019Standard struct {
}

// BuildId implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	return ""
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return JWK
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/secp256k1-2019/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) GetSupportedCurves() []string {
	return []string{"secp256k1"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return JWK | BASE58 | MULTIBASE | HEX
}

// GetType implements [IVerificationMethodStandard].
func (e *EcdsaSecp256k1VerificationKey2019Standard) GetType() string {
	return "EcdsaSecp256k1VerificationKey2019"
}

var _ IVerificationMethodStandard = (*Ed25519VerificationKey2018Standard)(nil)

type Ed25519VerificationKey2018Standard struct {
	multicodecSerializer IMulticodecSerializer
}

func NewEd25519VerificationKey2018Standard(multicodecSerializer IMulticodecSerializer) *Ed25519VerificationKey2018Standard {
	return &Ed25519VerificationKey2018Standard{multicodecSerializer: multicodecSerializer}
}

// BuildId implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	r, _ := e.multicodecSerializer.SerializePublicKey(asymmKey)
	return r
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return BASE58
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/ed25519-2018/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) GetSupportedCurves() []string {
	return []string{"Ed25519"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return BASE58
}

// GetType implements [IVerificationMethodStandard].
func (e *Ed25519VerificationKey2018Standard) GetType() string {
	return "Ed25519VerificationKey2018"
}

var _ IVerificationMethodStandard = (*JsonWebKey2020Standard)(nil)

type JsonWebKey2020Standard struct {
}

// BuildId implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	publicJWK, err := asymmKey.GetPublicJwk()
	if err != nil {
		return ""
	}

	if asymmKey.GetKty() == "OKP" {
		s, err := ComputeOKPThumbprint(publicJWK)
		if err != nil {
			return ""
		}
		return ToHex(s, false)
	}

	s, err := ComputeJWKThumbprint(publicJWK)
	if err != nil {
		return ""
	}
	return ToHex(s, false)
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return JWK
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/jws-2020/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) GetSupportedCurves() []string {
	return []string{}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return JWK
}

// GetType implements [IVerificationMethodStandard].
func (j *JsonWebKey2020Standard) GetType() string {
	return "JsonWebKey2020"
}

var _ IVerificationMethodStandard = (*X25519KeyAgreementKey2019Standard)(nil)

type X25519KeyAgreementKey2019Standard struct {
	multicodecSerializer IMulticodecSerializer
}

func NewX25519KeyAgreementKey2019Standard(multicodecSerializer IMulticodecSerializer) *X25519KeyAgreementKey2019Standard {
	return &X25519KeyAgreementKey2019Standard{multicodecSerializer: multicodecSerializer}
}

// BuildId implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	r, _ := x.multicodecSerializer.SerializePublicKey(asymmKey)
	return r
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return BASE58
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/x25519-2019/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) GetSupportedCurves() []string {
	return []string{"X25519"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return BASE58
}

// GetType implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2019Standard) GetType() string {
	return "X25519KeyAgreementKey2019"
}

var _ IVerificationMethodStandard = (*X25519KeyAgreementKey2020Standard)(nil)

type X25519KeyAgreementKey2020Standard struct {
}

// BuildId implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) BuildId(verificationMethod DidDocumentVerificationMethod, asymmKey IAsymmetricKey) string {
	return verificationMethod.PublicKeyMultibase
}

// GetDefaultEncoding implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) GetDefaultEncoding() SignatureKeyEncodingTypes {
	return MULTIBASE
}

// GetJSONLDContext implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) GetJSONLDContext() string {
	return "https://w3id.org/security/suites/x25519-2020/v1"
}

// GetSupportedCurves implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) GetSupportedCurves() []string {
	return []string{"X25519"}
}

// GetSupportedEncoding implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) GetSupportedEncoding() SignatureKeyEncodingTypes {
	return MULTIBASE
}

// GetType implements [IVerificationMethodStandard].
func (x *X25519KeyAgreementKey2020Standard) GetType() string {
	return "X25519KeyAgreementKey2020"
}
