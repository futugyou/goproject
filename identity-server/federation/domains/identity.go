package domains

type IdentityProviderMetadata struct {
	JwksUri                  string `json:"jwks_uri"`
	FastFedHandshakeStartUri string `json:"fastfed_handshake_start_uri"`
	BaseProviderMetadata
}
