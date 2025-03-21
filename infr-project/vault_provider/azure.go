package vault_provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type AzureClient struct {
	client *azsecrets.Client
}

func NewAzureClient() (*AzureClient, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	client, err := azsecrets.NewClient(os.Getenv("AZURE_VAULT_URL"), cred, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &AzureClient{
		client: client,
	}, nil
}

func (s *AzureClient) Search(ctx context.Context, key string) (*ProviderVault, error) {
	resp, err := s.client.GetSecret(ctx, key, "", nil)
	if err != nil {
		return nil, err
	}
	if resp.Value != nil {
		return &ProviderVault{
			Key:       resp.ID.Name(),
			Value:     *resp.Value,
			CreatedAt: *resp.Attributes.Created,
		}, nil
	}
	return nil, fmt.Errorf("can not found secret with name %s in Azure", key)
}

func (s *AzureClient) PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error) {
	var providerVaults = map[string]ProviderVault{}
	var keys = []string{}
	pager := s.client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return providerVaults, err
		}
		for _, secret := range page.Value {
			if strings.HasPrefix(secret.ID.Name(), prefix) {
				keys = append(keys, secret.ID.Name())
			}
		}
	}

	s.concurrencySearch(ctx, keys, providerVaults)

	return providerVaults, nil
}

func (s *AzureClient) concurrencySearch(ctx context.Context, keys []string, providerVaults map[string]ProviderVault) {
	var wg sync.WaitGroup
	concurrencyLimit := 5
	sem := make(chan struct{}, concurrencyLimit)
	errCh := make(chan error, len(keys))
	defer close(errCh)

	for _, key := range keys {
		wg.Add(1)

		go func(key string) {
			defer wg.Done()

			sem <- struct{}{}
			v, err := s.Search(ctx, key)
			if err != nil {
				errCh <- err
			} else {
				providerVaults[v.Key] = *v
			}

			<-sem
		}(key)
	}
}

func (s *AzureClient) BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error) {
	var providerVaults = map[string]ProviderVault{}
	var awskeys = []string{}
	pager := s.client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return providerVaults, err
		}
		for _, secret := range page.Value {
			for i := range keys {
				if secret.ID.Name() == keys[i] {
					awskeys = append(awskeys, secret.ID.Name())
				}
			}
		}
	}

	s.concurrencySearch(ctx, awskeys, providerVaults)

	return providerVaults, nil
}

func (s *AzureClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	resp, err := s.client.SetSecret(ctx, key, azsecrets.SetSecretParameters{Value: &value}, nil)
	if err != nil {
		return nil, err
	}
	return &ProviderVault{
		Key:       resp.ID.Name(),
		Value:     *resp.Value,
		CreatedAt: *resp.Attributes.Created,
	}, nil
}

func (s *AzureClient) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteSecret(ctx, key, nil)
	return err
}
