package did

import "time"

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
