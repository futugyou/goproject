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
