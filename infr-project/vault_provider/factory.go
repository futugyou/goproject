package vault_provider

import "fmt"

func VaultProviderFactory(provider string) (IVaultProviderAsync, error) {
	if provider == "AWS" {
		client, err := NewAWSClient()
		return NewAsyncWrapper(client), err
	}
	if provider == "HCP" {
		client, err := NewVaultClient()
		return NewAsyncWrapper(client), err
	}
	if provider == "AzureVault" {
		client, err := NewAzureClient()
		return NewAsyncWrapper(client), err
	}
	return nil, fmt.Errorf("provider type %s is not support", provider)
}
