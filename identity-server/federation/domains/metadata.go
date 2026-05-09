package domains

type ProviderMetadata struct {
	IdentityProvider    *IdentityProviderMetadata    `json:"identity_provider,omitempty"`
	ApplicationProvider *ApplicationProviderMetadata `json:"application_provider,omitempty"`
}
