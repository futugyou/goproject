package did

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/lestrrat-go/jwx/v4/jwa"
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

type DataIntegrityProof struct {
	ID                 string     `json:"id"`
	Type               string     `json:"type"`
	Created            *time.Time `json:"created,omitempty"`
	VerificationMethod string     `json:"verificationMethod"`
	ProofPurpose       string     `json:"proofPurpose"`
	Expires            *time.Time `json:"expires,omitempty"`
	Domain             []string   `json:"domain"`
	Challenge          string     `json:"challenge"`
	ProofValue         string     `json:"proofValue"`
	PreviousProof      []string   `json:"previousProof"`
	Nonce              string     `json:"nonce"`
	Jws                string     `json:"jws"`
}

type SecuredDocument struct {
	proofs                     []ISignatureProof
	verificationMethodEncoding IVerificationMethodEncoding
}

func NewSecuredDocument() *SecuredDocument {
	methods := []IVerificationMethod{
		&Ed25519VerificationMethod{},
		&Es256KVerificationMethod{},
		&Es256VerificationMethod{},
		&Es384VerificationMethod{},
		&X25519VerificationMethod{},
		&RSAVerificationMethod{},
		&JwkJcsPubVerificationMethod{},
	}
	return &SecuredDocument{
		proofs: []ISignatureProof{
			&Ed25519Signature2020Proof{},
			&JsonWebSignature2020Proof{},
		},
		verificationMethodEncoding: &VerificationMethodEncoding{
			standards:            GetAllVerificationMethodStandards(),
			multicodecSerializer: NewMulticodecSerializer(methods),
			verificationMethods:  methods,
		},
	}
}

func (s *SecuredDocument) Secure(document *W3CVerifiableCredential, didDocument *DidDocument, verificationMethodId string, asymKey IAsymmetricKey, createAt *time.Time) error {
	if document == nil || didDocument == nil || len(verificationMethodId) == 0 {
		return fmt.Errorf("parameter can not be nil")
	}

	var verificationMethod *DidDocumentVerificationMethod
	for _, v := range didDocument.VerificationMethod {
		if v.ID == verificationMethodId {
			verificationMethod = &v
			break
		}
	}

	if verificationMethod == nil {
		return fmt.Errorf("verification method not found in did document")
	}

	var proof ISignatureProof
	for _, v := range s.proofs {
		if v.GetVerificationMethod() == verificationMethod.Type {
			proof = v
		}
	}

	if proof == nil {
		return fmt.Errorf("proof not found for verification method")
	}

	created := time.Now().UTC()
	if createAt != nil {
		created = *createAt
	}

	dataIntegrityProof := &DataIntegrityProof{
		Type:               proof.GetType(),
		Created:            &created,
		VerificationMethod: verificationMethod.ID,
		ProofPurpose:       "assertionMethod",
	}

	docBytes, _ := json.Marshal(document)
	var docMap map[string]any
	json.Unmarshal(docBytes, &docMap)

	hashPayload, _ := s.hashDocument(docMap)
	hashProof, _ := s.hashProof(dataIntegrityProof, docMap)

	result := append(hashProof, hashPayload...)

	if asymKey == nil {
		asymKey, _ = s.verificationMethodEncoding.Decode(*verificationMethod)
	}

	err := proof.ComputeProof(dataIntegrityProof, result, asymKey, proof.GetHashingMethod())
	if err != nil {
		return err
	}

	proofBytes, err := json.Marshal(dataIntegrityProof)
	if err != nil {
		return fmt.Errorf("failed to marshal proof: %w", err)
	}
	document.Proof = json.RawMessage(proofBytes)

	return nil
}

func (s *SecuredDocument) hash(jsonData any) ([]byte, error) {
	canonized, err := RdfTransform(jsonData)
	if err != nil {
		return nil, fmt.Errorf("canonization failed: %w", err)
	}

	hash := sha256.Sum256([]byte(canonized))
	return hash[:], nil
}

func (s *SecuredDocument) hashDocument(docMap map[string]any) ([]byte, error) {
	delete(docMap, "proof")
	return s.hash(docMap)
}

func (s *SecuredDocument) hashProof(diProof *DataIntegrityProof, docMap map[string]any) ([]byte, error) {
	proofBytes, _ := json.Marshal(diProof)
	var pMap map[string]any
	json.Unmarshal(proofBytes, &pMap)

	if ctx, ok := docMap["@context"]; ok {
		pMap["@context"] = ctx
	}

	return s.hash(pMap)
}

func (s *SecuredDocument) getProof(pt string) ISignatureProof {
	for _, v := range s.proofs {
		if v.GetVerificationMethod() == pt {
			return v
		}
	}

	return nil
}

func (s *SecuredDocument) Check(verifiableDocument *W3CVerifiableCredential, didDocument *DidDocument) (bool, error) {
	if len(verifiableDocument.Proof) == 0 {
		return false, errors.New("the document doesn't contain a proof")
	}

	var diProof DataIntegrityProof
	if err := json.Unmarshal(verifiableDocument.Proof, &diProof); err != nil {
		return false, fmt.Errorf("failed to unmarshal proof: %w", err)
	}

	vm := didDocument.GetVerificationMethod(diProof.VerificationMethod)
	if vm == nil {
		return false, fmt.Errorf("verification method not found in DID document")
	}

	docBytes, _ := json.Marshal(verifiableDocument)
	var docMap map[string]any
	json.Unmarshal(docBytes, &docMap)

	proof := s.getProof(vm.Type)

	emptyProof := diProof
	emptyProof.ProofValue = ""
	emptyProof.Jws = ""

	hashPayload, _ := s.hashDocument(docMap)
	hashProof, _ := s.hashProof(&emptyProof, docMap)

	result := append(hashProof, hashPayload...)

	signature, _ := proof.GetSignature(&diProof)
	asymKey, _ := s.verificationMethodEncoding.Decode(*vm)

	return asymKey.CheckHash(result, signature, jwa.NewSignatureAlgorithm(proof.GetHashingMethod())), nil
}

func (s *SecuredDocument) SecureJwt(
	subject string,
	didDocument *DidDocument,
	verificationMethodId string,
	vcCredential *W3CVerifiableCredential,
	asymKey IAsymmetricKey,
) (string, error) {
	vm := didDocument.GetVerificationMethod(verificationMethodId)
	if vm == nil {
		return "", fmt.Errorf("verification method %s not found", verificationMethodId)
	}

	if asymKey == nil {
		if vm.PrivateKeyJwk != nil {
			k, _ := jwk.ParseKey(vm.PrivateKeyJwk)
			asymKey = NewJsonWebKeySecurityKey(k)
		} else {
			asymKey, _ = s.verificationMethodEncoding.Decode(*vm)
		}
	}

	if asymKey == nil {
		return "", fmt.Errorf("failed to obtain asymmetric key for signing")
	}

	vcBytes, _ := json.Marshal(vcCredential)
	var vcMap map[string]any
	json.Unmarshal(vcBytes, &vcMap)

	token := jwt.New()
	token.Set("sub", subject)
	token.Set("jti", vcCredential.Id)
	token.Set("vc", vcMap)
	token.Set("iss", didDocument.Id)

	now := time.Now().UTC()
	token.Set("iat", now.Unix())
	token.Set("nbf", now.Unix())

	privateJwk, err := asymKey.GetPrivateJwk()
	if err != nil {
		return "", err
	}

	alg := jwa.NewSignatureAlgorithm(asymKey.GetJwtAlg())

	signedBytes, err := jwt.Sign(token, jwt.WithKey(alg, privateJwk))
	if err != nil {
		return "", fmt.Errorf("impossible to create a JWT: %w", err)
	}

	return string(signedBytes), nil
}
