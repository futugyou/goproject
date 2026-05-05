package did

import (
	"encoding/json"
	"time"

	"github.com/lestrrat-go/jwx/v4/jwk"
)

type DecentralizedIdentifier struct {
	Scheme     string
	Method     string
	Identifier string
	Fragment   string
}

func (d DecentralizedIdentifier) GetDidWithoutFragment() string {
	return d.Scheme + ":" + d.Method + ":" + d.Identifier
}

func NewDecentralizedIdentifier(scheme, method, identifier, fragment string) DecentralizedIdentifier {
	return DecentralizedIdentifier{
		Scheme:     scheme,
		Method:     method,
		Identifier: identifier,
		Fragment:   fragment,
	}
}

type DidDocumentService struct {
	ID              string `json:"id"`
	ServiceEndpoint string `json:"serviceEndpoint"`
	Type            string `json:"type"`
}

type VerificationMethodUsage uint32

const (
	Authentication       VerificationMethodUsage = 1 << iota // 1 (2^0)
	AssertionMethod                                          // 2 (2^1)
	KeyAgreement                                             // 4 (2^2)
	CapabilityInvocation                                     // 8 (2^3)
	CapabilityDelegation                                     // 16 (2^4)
)

type KeyPurposes uint32

const (
	SigAuthentication VerificationMethodUsage = 0
	VerificationKey   VerificationMethodUsage = 1
	Encryption        VerificationMethodUsage = 2
)

type DidDocumentVerificationMethod struct {
	ID                   string                  `json:"id"`
	Type                 string                  `json:"type"`
	Controller           string                  `json:"controller"`
	PublicKeyJwk         json.RawMessage         `json:"publicKeyJwk"`
	PrivateKeyJwk        json.RawMessage         `json:"privateKeyJwk"`
	PublicKeyMultibase   string                  `json:"publicKeyMultibase"`
	SecretKeyMultibase   string                  `json:"secretKeyMultibase"`
	BlockChainAccountId  string                  `json:"blockchainAccountId"`
	PublicKeyHex         string                  `json:"publicKeyHex"`
	PrivateKeyHex        string                  `json:"privateKeyHex"`
	PublicKeyBase64      string                  `json:"publicKeyBase64"`
	PublicKeyBase58      string                  `json:"publicKeyBase58"`
	PrivateKeyBase58     string                  `json:"privateKeyBase58"`
	PublicKeyPem         string                  `json:"publicKeyPem"`
	Value                string                  `json:"value"`
	Usage                VerificationMethodUsage `json:"-"`
	AdditionalParameters map[string]string       `json:"-"`
	JsonLdContext        string                  `json:"-"`
}

func (v *DidDocumentVerificationMethod) GetPublicKeyJwk() (jwk.Key, error) {
	return jwk.ParseKey(v.PublicKeyJwk)
}

func (v *DidDocumentVerificationMethod) GetPrivateKeyJwk() (jwk.Key, error) {
	return jwk.ParseKey(v.PrivateKeyJwk)
}

type IdentityDocumentIdentifier struct {
	Identifier string
	Source     string
	Address    string
	PublicKey  string
}

type DidDocument struct {
	Context              any                             `json:"@context"`
	Id                   string                          `json:"id"`
	Controller           any                             `json:"controller,omitempty"`
	AlsoKnownAs          []string                        `json:"alsoKnownAs,omitempty"`
	VerificationMethod   []DidDocumentVerificationMethod `json:"verificationMethod,omitempty"`
	Authentication       []any                           `json:"authentication,omitempty"`
	AssertionMethod      []any                           `json:"assertionMethod,omitempty"`
	KeyAgreement         []any                           `json:"keyAgreement,omitempty"`
	CapabilityInvocation []any                           `json:"capabilityInvocation,omitempty"`
	CapabilityDelegation []any                           `json:"capabilityDelegation,omitempty"`
	Service              []DidDocumentService            `json:"service,omitempty"`
}

func (d *DidDocument) GetVerificationMethod(methodID string) *DidDocumentVerificationMethod {
	for _, v := range d.VerificationMethod {
		if v.ID == methodID {
			return &v
		}
	}

	return nil
}

func (d *DidDocument) Serialize() (string, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (d *DidDocument) AddVerificationMethod(verificationMethod DidDocumentVerificationMethod, isStandard bool) {
	d.VerificationMethod = append(d.VerificationMethod, verificationMethod)
}

func (d *DidDocument) AddService(service DidDocumentService, isStandard bool) {
	d.Service = append(d.Service, service)
}

type W3CVerifiableCredential struct {
	Proof             json.RawMessage               `json:"proof"`
	Context           []string                      `json:"@context"`
	Id                string                        `json:"id"`
	Types             []string                      `json:"type"`
	Issuer            string                        `json:"issuer"`
	IssuanceDate      time.Time                     `json:"issuanceDate"`
	ValidFrom         *time.Time                    `json:"validFrom,omitempty"`
	ValidUntil        *time.Time                    `json:"validUntil,omitempty"`
	Issued            time.Time                     `json:"issued"`
	CredentialSubject any                           `json:"credentialSubject"`
	Schema            W3CVerifiableCredentialSchema `json:"credentialSchema"`
}

func (vc *W3CVerifiableCredential) AddSubjectRecord(record map[string]any) error {
	if vc.CredentialSubject == nil {
		vc.CredentialSubject = record
		return nil
	}

	switch existing := vc.CredentialSubject.(type) {
	case map[string]any:
		vc.CredentialSubject = []any{existing, record}
	case []any:
		vc.CredentialSubject = append(existing, record)
	case json.RawMessage:
		var temp any
		if err := json.Unmarshal(existing, &temp); err != nil {
			return err
		}
		vc.CredentialSubject = temp
		return vc.AddSubjectRecord(record)
	default:
		vc.CredentialSubject = []any{record}
	}

	return nil
}

type W3CVerifiableCredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
