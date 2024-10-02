package vault_provider

import (
	"context"
	"fmt"
	"time"
)

type AzureClient struct {
}

func NewAzureClient() *AzureClient {
	return &AzureClient{}
}

func (s *AzureClient) Get(ctx context.Context, key string) (*ProviderVault, error) {
	return nil, fmt.Errorf("can not found secret with name %s in Azure", key)
}

func (s *AzureClient) Search(ctx context.Context, prefix string) ([]ProviderVault, error) {
	var providerVaults = []ProviderVault{}

	return providerVaults, nil
}

func (s *AzureClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	return &ProviderVault{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *AzureClient) Delete(ctx context.Context, key string) error {
	return nil
}
