package vault_provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	vault "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-06-13/client/secret_service"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
)

type VaultClient struct {
	http vault.ClientService
}

func NewVaultClient() (*VaultClient, error) {
	hcpConfig, err := config.NewHCPConfig(
		config.FromEnv(),
	)
	if err != nil {
		return nil, err
	}

	// Construct HTTP client config
	httpclientConfig := httpclient.Config{
		HCPConfig: hcpConfig,
	}

	// Initialize SDK http client
	cl, err := httpclient.New(httpclientConfig)
	if err != nil {
		return nil, err
	}

	// Import versioned client for each service.
	vaultClient := vault.New(cl, nil)
	return &VaultClient{
		http: vaultClient,
	}, err
}

func (s *VaultClient) Search(ctx context.Context, key string) (*ProviderVault, error) {
	params := &vault.OpenAppSecretParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		SecretName:             key,
		Context:                ctx,
	}

	var result *vault.OpenAppSecretOK
	var err error

	if result, err = s.http.OpenAppSecret(params, nil); err != nil {
		return nil, err
	}

	if result.Payload != nil && result.Payload.Secret != nil && result.Payload.Secret.Version != nil {
		return &ProviderVault{
			Key:       key,
			Value:     result.Payload.Secret.Version.Value,
			CreatedAt: time.Time(result.Payload.Secret.Version.CreatedAt),
		}, nil
	}

	return nil, fmt.Errorf("can not found secret with name %s in hashicorp", key)
}

func (s *VaultClient) PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error) {
	params := &vault.OpenAppSecretsParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		Context:                ctx,
	}
	var result *vault.OpenAppSecretsOK
	var err error
	if result, err = s.http.OpenAppSecrets(params, nil); err != nil {
		return nil, err
	}

	var providerVaults = map[string]ProviderVault{}
	if result.Payload != nil && len(result.Payload.Secrets) > 0 {
		for _, v := range result.Payload.Secrets {
			if v != nil && v.Version != nil {
				if len(prefix) > 0 {
					if strings.HasPrefix(v.Version.Value, prefix) {
						providerVaults[v.Name] = ProviderVault{
							Key:       v.Name,
							Value:     v.Version.Value,
							CreatedAt: time.Time(v.Version.CreatedAt),
						}
					}
				} else {
					providerVaults[v.Name] = ProviderVault{
						Key:       v.Name,
						Value:     v.Version.Value,
						CreatedAt: time.Time(v.Version.CreatedAt),
					}
				}
			}
		}
		return providerVaults, nil
	}

	return nil, fmt.Errorf("can not found secret with prefix %s in hashicorp", prefix)
}

func (s *VaultClient) BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error) {
	params := &vault.OpenAppSecretsParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		Context:                ctx,
	}
	var result *vault.OpenAppSecretsOK
	var err error
	if result, err = s.http.OpenAppSecrets(params, nil); err != nil {
		return nil, err
	}

	var providerVaults = map[string]ProviderVault{}
	if result.Payload != nil && len(result.Payload.Secrets) > 0 {
		for _, v := range result.Payload.Secrets {
			if v != nil && v.Version != nil {
				for i := range keys {
					if keys[i] == v.Name {
						providerVaults[v.Name] = ProviderVault{
							Key:       v.Name,
							Value:     v.Version.Value,
							CreatedAt: time.Time(v.Version.CreatedAt),
						}
					}
				}
			}
		}
	}

	return providerVaults, nil
}

func (s *VaultClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	params := &vault.CreateAppKVSecretParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		Context:                ctx,
		Body: vault.CreateAppKVSecretBody{
			Name:  key,
			Value: value,
		},
	}

	var result *vault.CreateAppKVSecretOK
	var err error
	if result, err = s.http.CreateAppKVSecret(params, nil); err != nil {
		return nil, err
	}

	if result.Payload != nil || result.Payload.Secret != nil && result.Payload.Secret.Version != nil {
		return &ProviderVault{
			Key:       key,
			Value:     value,
			CreatedAt: time.Time(result.Payload.Secret.Version.CreatedAt),
		}, nil
	}

	return nil, fmt.Errorf("call hashicorp ok, but same thing error, try to find detail in hashicorp cloud")
}

func (s *VaultClient) Delete(ctx context.Context, key string) error {
	params := &vault.DeleteAppSecretParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		Context:                ctx,
		SecretName:             key,
	}

	_, err := s.http.DeleteAppSecret(params, nil)
	return err
}

func (s *VaultClient) SearchAsync(ctx context.Context, key string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Search(ctx, key)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *VaultClient) PrefixSearchAsync(ctx context.Context, prefix string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.PrefixSearch(ctx, prefix)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *VaultClient) BatchSearchAsync(ctx context.Context, keys []string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.BatchSearch(ctx, keys)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *VaultClient) UpsertAsync(ctx context.Context, key string, value string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Upsert(ctx, key, value)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *VaultClient) DeleteAsync(ctx context.Context, key string) <-chan error {
	errorChan := make(chan error, 1)
	go func() {
		defer close(errorChan)
		errorChan <- s.Delete(ctx, key)
	}()
	return errorChan
}
