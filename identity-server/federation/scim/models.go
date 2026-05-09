package scim

const JwtProfile string = "urn:ietf:params:fastfed:1.0:provider_authentication:oauth:2.0:jwt_profile"
const SchemaGrammarName string = "urn:ietf:params:fastfed:1.0:schemas:scim:2.0"
const ProvisioningProfileName string = "urn:ietf:params:fastfed:1.0:provisioning:scim:2.0:enterprise"

type AuthenticationProfileResult struct {
	TokenEndpoint string `json:"token_endpoint"`
	Scope         string `json:"scope"`
}

type ScimEntrepriseMappingsResult struct {
	CanSupportedNestedGroups  bool              `json:"can_support_nested_groups"`
	MaxGroupMembershipChanges int               `json:"max_group_membership_changes"`
	DesiredAttributes         DesiredAttributes `json:"desired_attributes"`
}

type DesiredAttributes struct {
	Attrs *SchemaGrammarDesiredAttributes `json:"urn:ietf:params:fastfed:1.0:schemas:scim:2.0"`
}

type SchemaGrammarDesiredAttributes struct {
	RequiredUserAttributes  []string `json:"required_user_attributes"`
	OptionalUserAttributes  []string `json:"optional_user_attributes"`
	RequiredGroupAttributes []string `json:"required_group_attributes"`
	OptionalGroupAttributes []string `json:"optional_group_attributes"`
}

type ScimEntrepriseRegistrationResult struct {
	ScimServiceUri                string                      `json:"scim_service_uri"`
	ProviderAuthenticationMethods string                      `json:"provider_authentication_methods"`
	JwtProfile                    AuthenticationProfileResult `json:"urn:ietf:params:fastfed:1.0:provider_authentication:oauth:2.0:jwt_profile"`
}
