package provider

import (
	"fmt"

	"github.com/futugyou/vaultservice/options"
)

func VaultProviderFactory(provider string, opts *options.Options) (VaultProvider, error) {
	if opts == nil {
		return nil, fmt.Errorf("opts is nil")
	}

	if provider == "AWS" {
		return NewAWSClient(opts)
	}

	if provider == "HCP" {
		return NewVaultClient(opts)
	}

	if provider == "AzureVault" {
		return NewAzureClient(opts)
	}

	return nil, fmt.Errorf("provider type %s is not support", provider)
}
