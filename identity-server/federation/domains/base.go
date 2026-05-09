package domains

import (
	"fmt"
)

type ProviderContactInformation struct {
	Organization string `json:"organization"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
}

func (p ProviderContactInformation) Validate() []string {
	result := []string{}
	if len(p.Organization) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "organization"))
	}
	if len(p.Email) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "email"))
	}
	return result
}

type DisplaySettings struct {
	DisplayName string `json:"display_name"`
	IconUri     string `json:"icon_uri"`
	LogoUri     string `json:"logo_uri"`
	License     string `json:"license"`
}

func (p DisplaySettings) Validate() []string {
	result := []string{}
	if len(p.DisplayName) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "display_name"))
	}
	if len(p.License) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "license"))
	}
	return result
}

type Capabilities struct {
	AuthenticationProfiles []string `json:"authentication_profiles"`
	ProvisioningProfiles   []string `json:"provisioning_profiles"`
	SchemaGrammars         []string `json:"schema_grammars"`
	SigningAlgorithms      []string `json:"signing_algorithms"`
}

func (p Capabilities) Validate() []string {
	result := []string{}
	if len(p.SigningAlgorithms) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "signing_algorithms"))
	}
	if len(p.SchemaGrammars) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "schema_grammars"))
	}
	return result
}

type BaseProviderMetadata struct {
	EntityId           string                     `json:"entity_id"`
	ProviderDomain     string                     `json:"provider_domain"`
	ContactInformation ProviderContactInformation `json:"provider_contact_information"`
	DisplaySettings    DisplaySettings            `json:"display_settings"`
	Capabilities       Capabilities               `json:"capabilities"`
}

func (p *BaseProviderMetadata) Validate() []string {
	result := []string{}
	if len(p.EntityId) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "entity_id"))
	}
	if len(p.ProviderDomain) == 0 {
		result = append(result, fmt.Sprintf("Parameter '%s' is missing from the metadata", "provider_domain"))
	}
	result = append(result, p.ContactInformation.Validate()...)
	result = append(result, p.DisplaySettings.Validate()...)
	result = append(result, p.Capabilities.Validate()...)
	return result
}

func (p BaseProviderMetadata) CheckCompatibility(cap Capabilities) []string {
	var result []string

	check := func(src, target []string, format string) {
		if src == nil || target == nil {
			return
		}

		targetMap := make(map[string]struct{})
		for _, v := range target {
			targetMap[v] = struct{}{}
		}

		for _, v := range src {
			if _, exists := targetMap[v]; !exists {
				result = append(result, fmt.Sprintf(format, v))
			}
		}
	}

	check(cap.AuthenticationProfiles, p.Capabilities.AuthenticationProfiles, "Authentication profile %s is not compatible")
	check(cap.ProvisioningProfiles, p.Capabilities.ProvisioningProfiles, "Provisioning profile %s is not compatible")
	check(cap.SchemaGrammars, p.Capabilities.SchemaGrammars, "Schema Grammar %s is not compatible")
	check(cap.SigningAlgorithms, p.Capabilities.SigningAlgorithms, "Signing Algorithm %s is not compatible")

	return result
}
