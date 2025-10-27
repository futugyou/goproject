package vault_provider

import "fmt"

func VaultProviderFactory(provider string) (IVaultProvider, error) {
	if provider == "AWS" {
		return NewAWSClient()
	}
	if provider == "HCP" {
		return NewVaultClient()
	}
	if provider == "AzureVault" {
		return NewAzureClient()
	}
	return nil, fmt.Errorf("provider type %s is not support", provider)
}
