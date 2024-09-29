package vault_provider

import (
	"context"
	"fmt"
	"log"
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

func NewVaultClient() *VaultClient {
	hcpConfig, err := config.NewHCPConfig(
		config.FromEnv(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Construct HTTP client config
	httpclientConfig := httpclient.Config{
		HCPConfig: hcpConfig,
	}

	// Initialize SDK http client
	cl, err := httpclient.New(httpclientConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Import versioned client for each service.
	vaultClient := vault.New(cl, nil)
	return &VaultClient{
		http: vaultClient,
	}
}

func (s *VaultClient) Get(ctx context.Context, key string) (*ProviderVault, error) {
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

	if result.Payload != nil || result.Payload.Secret != nil && result.Payload.Secret.Version != nil {
		return &ProviderVault{
			Key:       key,
			Value:     result.Payload.Secret.Version.Value,
			CreatedAt: time.Time(result.Payload.Secret.Version.CreatedAt),
		}, nil
	}

	return nil, fmt.Errorf("can not found secret with name %s in hashicorp", key)
}

func (s *VaultClient) Search(ctx context.Context, key string) ([]ProviderVault, error) {
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

	var providerVaults = []ProviderVault{}
	if result.Payload != nil || len(result.Payload.Secrets) == 0 {
		for _, v := range result.Payload.Secrets {
			if v != nil && v.Version != nil {
				if len(key) > 0 {
					if strings.HasPrefix(v.Version.Value, key) {
						providerVaults = append(providerVaults, ProviderVault{
							Key:       key,
							Value:     v.Version.Value,
							CreatedAt: time.Time(v.Version.CreatedAt),
						})
					}
				} else {
					providerVaults = append(providerVaults, ProviderVault{
						Key:       key,
						Value:     v.Version.Value,
						CreatedAt: time.Time(v.Version.CreatedAt),
					})
				}
			}
		}
		return providerVaults, nil
	}

	return nil, fmt.Errorf("can not found secret with name %s in hashicorp", key)
}

func (s *VaultClient) Upinsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	return nil, nil
}

func (s *VaultClient) Delete(ctx context.Context, key string) error {

	return nil
}
