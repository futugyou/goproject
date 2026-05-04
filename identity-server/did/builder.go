package did

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"
)

// builder := NewContextBuilder()
//
//	result := builder.
//		AddPropertyValue("name", "http://schema.org/name").
//		AddPropertyId("owner", "http://schema.org/owner").
//		Build()
//
//	jsonBytes, _ := json.MarshalIndent(result, "", "  ")
//	fmt.Println(string(jsonBytes))
type ContextBuilder struct {
	jObj map[string]any
}

func NewContextBuilder() *ContextBuilder {
	return &ContextBuilder{
		jObj: make(map[string]any),
	}
}

func (b *ContextBuilder) AddPropertyValue(shortName, fullPath string) *ContextBuilder {
	b.jObj[shortName] = fullPath
	return b
}

func (b *ContextBuilder) AddPropertyId(shortName, fullPath string) *ContextBuilder {
	b.jObj[shortName] = map[string]string{
		"@id":   fullPath,
		"@type": "@id",
	}
	return b
}

func (b *ContextBuilder) AddRaw(shortName string, raw []byte) *ContextBuilder {
	b.jObj[shortName] = json.RawMessage(raw)
	return b
}

func (b *ContextBuilder) Build() map[string]any {
	return b.jObj
}

const DIDContext string = "https://www.w3.org/ns/did/v1"

type DidDocumentBuilder struct {
	contextUrls                []string
	context                    []any
	controllers                []string
	innerVerificationMethods   []DidDocumentVerificationMethod
	identityDocument           *DidDocument
	verificationMethodEncoding IVerificationMethodEncoding
}

func NewDidDocumentBuilder(identityDocument DidDocument, verificationMethodEncoding IVerificationMethodEncoding) *DidDocumentBuilder {
	return &DidDocumentBuilder{
		identityDocument:           &identityDocument,
		verificationMethodEncoding: verificationMethodEncoding,
	}
}

func (b *DidDocumentBuilder) AddContext(callback func(*ContextBuilder)) *DidDocumentBuilder {
	builder := &ContextBuilder{}
	callback(builder)
	b.context = append(b.context, builder.Build())
	return b
}

func (b *DidDocumentBuilder) AddContextUrls(contextList []string) *DidDocumentBuilder {
	for _, v := range contextList {
		index := slices.IndexFunc(b.contextUrls, func(url string) bool {
			return v == url
		})
		if index == -1 {
			b.contextUrls = append(b.contextUrls, v)
		}
	}
	return b
}

func (b *DidDocumentBuilder) AddController(controller string) *DidDocumentBuilder {
	b.controllers = append(b.controllers, controller)
	return b
}

func (b *DidDocumentBuilder) AddAlsoKnownAs(alsoKnownAs string) *DidDocumentBuilder {
	if b.identityDocument == nil {
		b.identityDocument = &DidDocument{}
	}

	b.identityDocument.AlsoKnownAs = append(b.identityDocument.AlsoKnownAs, alsoKnownAs)
	return b
}

func (b *DidDocumentBuilder) AddVerificationMethod(verificationMethodStandard string,
	asymmKey IAsymmetricKey,
	controller string,
	usage VerificationMethodUsage,
	isReference bool,
	includePrivateKey bool,
	encoding SignatureKeyEncodingTypes,
	callback func(*DidDocumentVerificationMethod),
) *DidDocumentBuilder {
	verificationMethod, _ := b.verificationMethodEncoding.Encode(
		verificationMethodStandard,
		b.identityDocument.Id,
		controller,
		asymmKey,
		encoding,
		includePrivateKey)

	if verificationMethod == nil {
		return b
	}

	if callback != nil {
		callback(verificationMethod)
	}
	standards := b.verificationMethodEncoding.GetStandards()
	index := slices.IndexFunc(standards, func(stn IVerificationMethodStandard) bool {
		return stn.GetType() == verificationMethod.Type
	})

	standard := standards[index]

	jsonld := standard.GetJSONLDContext()
	if len(jsonld) > 0 {
		b.AddContextUrls([]string{jsonld})
	}

	verificationMethod.Usage = usage
	if len(verificationMethod.ID) == 0 {
		index := len(b.identityDocument.VerificationMethod) + len(b.innerVerificationMethods) + 1
		verificationMethod.ID = verificationMethod.Controller + "#keys-" + string(rune(index))
	}

	if isReference {
		b.identityDocument.AddVerificationMethod(*verificationMethod, true)
	} else {
		b.innerVerificationMethods = append(b.innerVerificationMethods, *verificationMethod)
	}

	return b
}

func (b *DidDocumentBuilder) AddServiceEndpoint(serviceType string, serviceEndpoint string) *DidDocumentBuilder {
	if b.identityDocument == nil {
		b.identityDocument = &DidDocument{}
	}

	var id = fmt.Sprintf("%s#service-%d", b.identityDocument.Id, len(b.identityDocument.Service)+1)
	b.identityDocument.AddService(DidDocumentService{
		ID:              id,
		Type:            serviceType,
		ServiceEndpoint: serviceEndpoint,
	}, false)
	return b
}

func (b *DidDocumentBuilder) AddDidDocumentService(service DidDocumentService) *DidDocumentBuilder {
	if b.identityDocument == nil {
		b.identityDocument = &DidDocument{}
	}

	b.identityDocument.Service = append(b.identityDocument.Service, service)
	return b
}

func (b *DidDocumentBuilder) Build() *DidDocument {
	if b.identityDocument == nil {
		b.identityDocument = &DidDocument{}
	}
	b.identityDocument.Context = b.buildContext()
	b.identityDocument.Controller = b.buildController()
	b.identityDocument.Authentication = b.buildEmbeddedVerificationMethods(Authentication)
	b.identityDocument.AssertionMethod = b.buildEmbeddedVerificationMethods(AssertionMethod)
	b.identityDocument.KeyAgreement = b.buildEmbeddedVerificationMethods(KeyAgreement)
	b.identityDocument.CapabilityInvocation = b.buildEmbeddedVerificationMethods(CapabilityInvocation)
	b.identityDocument.CapabilityDelegation = b.buildEmbeddedVerificationMethods(CapabilityDelegation)
	return b.identityDocument
}

func (b *DidDocumentBuilder) buildEmbeddedVerificationMethods(usage VerificationMethodUsage) []any {
	var referencedVerificationMethods []DidDocumentVerificationMethod
	var innerVerificationMethods []DidDocumentVerificationMethod
	for _, v := range b.identityDocument.VerificationMethod {
		if v.Usage|usage != 0 {
			referencedVerificationMethods = append(referencedVerificationMethods, v)
		}
	}
	for _, v := range b.innerVerificationMethods {
		if v.Usage|usage != 0 {
			innerVerificationMethods = append(innerVerificationMethods, v)
		}
	}
	if len(referencedVerificationMethods) == 0 && len(innerVerificationMethods) == 0 {
		return nil
	}

	result := make([]any, 0, len(referencedVerificationMethods)+len(innerVerificationMethods))

	for _, v := range referencedVerificationMethods {
		result = append(result, v)
	}

	for _, v := range innerVerificationMethods {
		result = append(result, v)
	}

	return result
}

func (b *DidDocumentBuilder) buildController() any {
	if len(b.controllers) == 0 {
		return nil
	}
	if len(b.controllers) == 01 {
		return b.controllers[0]
	}
	return b.controllers
}

func (b *DidDocumentBuilder) buildContext() any {
	if len(b.contextUrls) == 0 && len(b.context) == 0 {
		return DIDContext
	}
	result := []any{DIDContext}
	for _, v := range b.contextUrls {
		result = append(result, v)
	}
	for _, v := range b.context {
		result = append(result, v)
	}

	return result
}

// builder := NewCredentialSubjectBuilder("did:example:123")
// builder.SetPersonalIdentifier("PID-999")
// builder.AddClaim("age", 30)
// builder.AddClaim("address", map[string]any{
// "city": "Shanghai",
// })
// result := builder.Build()
// bytes, _ := json.MarshalIndent(result, "", "  ")
// fmt.Println(string(bytes))
type CredentialSubjectBuilder struct {
	jsonObj map[string]any
}

func NewCredentialSubjectBuilder(id string) *CredentialSubjectBuilder {
	obj := make(map[string]any)
	if id != "" {
		obj["id"] = id
	}
	return &CredentialSubjectBuilder{jsonObj: obj}
}

func (b *CredentialSubjectBuilder) SetPersonalIdentifier(personalIdentifier string) *CredentialSubjectBuilder {
	b.jsonObj["personalIdentifier"] = personalIdentifier
	return b
}

func (b *CredentialSubjectBuilder) AddClaim(name string, value any) *CredentialSubjectBuilder {
	b.jsonObj[name] = value
	return b
}

func (b *CredentialSubjectBuilder) Build() map[string]any {
	return b.jsonObj
}

type VcBuilder struct {
	credential *W3CVerifiableCredential
}

func NewVcBuilder(credential *W3CVerifiableCredential) *VcBuilder {
	return &VcBuilder{credential: credential}
}

func NewVcComplexBuilder(id,
	jsonLdContext,
	issuer,
	credentialType string,
	additionalTypes []string,
	validFrom *time.Time,
	validUntil *time.Time,
) *VcBuilder {
	now := time.Now().UTC()
	credential := &W3CVerifiableCredential{
		Id:           id,
		Issuer:       issuer,
		ValidFrom:    validFrom,
		ValidUntil:   validUntil,
		IssuanceDate: now,
		Issued:       now,
	}

	credential.Context = append(credential.Context, "https://www.w3.org/2018/credentials/v1")
	credential.Context = append(credential.Context, jsonLdContext)
	credential.Types = append(credential.Types, "VerifiableCredential")
	credential.Types = append(credential.Types, additionalTypes...)
	credential.Types = append(credential.Types, credentialType)

	return NewVcBuilder(credential)
}

func (b *VcBuilder) AddCredentialSubject(id string, action func(*CredentialSubjectBuilder)) *VcBuilder {
	builder := NewCredentialSubjectBuilder(id)
	action(builder)
	record := builder.Build()
	b.credential.AddSubjectRecord(record)
	return b
}

func (b *VcBuilder) SetSchema(id, schematype string) *VcBuilder {
	b.credential.Schema = W3CVerifiableCredentialSchema{
		ID:   id,
		Type: schematype,
	}

	return b
}

func (b *VcBuilder) Build() *W3CVerifiableCredential {
	return b.credential
}
