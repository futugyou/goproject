package vault_provider

import "fmt"

func VaultProviderFatory(provider string) (IVaultProvider, error) {
	if provider == "AWS" {
		return NewAWSClient(), nil
	}
	if provider == "HCP" {
		return NewVaultClient(), nil
	}
	if provider == "AzureVault" {
		return NewAzureClient(), nil
	}
	return nil, fmt.Errorf("provider type %s is not support", provider)
}
