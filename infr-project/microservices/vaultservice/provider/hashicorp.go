package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/futugyou/vaultservice/options"
	vault "github.com/hashicorp/hcp-sdk-go/clients/cloud-vault-secrets/stable/2023-06-13/client/secret_service"
	"github.com/hashicorp/hcp-sdk-go/config"
	"github.com/hashicorp/hcp-sdk-go/httpclient"
)

type hashicorpClient struct {
	http vault.ClientService
	opts *options.Options
}

func newVaultClient(opts *options.Options) (*hashicorpClient, error) {
	var configs []config.HCPConfigOption
	if len(opts.HcpClientID) > 0 && len(opts.HcpClientSecret) > 0 {
		configs = append(configs, config.WithClientCredentials(opts.HcpClientID, opts.HcpClientSecret))
	} else {
		configs = append(configs, config.FromEnv())
	}

	hcpConfig, err := config.NewHCPConfig(configs...)
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
	return &hashicorpClient{
		http: vaultClient,
		opts: opts,
	}, err
}

func (s *hashicorpClient) Search(ctx context.Context, key string) (*ProviderVault, error) {
	params := &vault.OpenAppSecretParams{
		AppName:                s.opts.HcpAppName,
		LocationOrganizationID: s.opts.HcpOrganizationID,
		LocationProjectID:      s.opts.HcpProjectID,
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

func (s *hashicorpClient) PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error) {
	params := &vault.OpenAppSecretsParams{
		AppName:                s.opts.HcpAppName,
		LocationOrganizationID: s.opts.HcpOrganizationID,
		LocationProjectID:      s.opts.HcpProjectID,
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

func (s *hashicorpClient) BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error) {
	params := &vault.OpenAppSecretsParams{
		AppName:                s.opts.HcpAppName,
		LocationOrganizationID: s.opts.HcpOrganizationID,
		LocationProjectID:      s.opts.HcpProjectID,
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

func (s *hashicorpClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	params := &vault.CreateAppKVSecretParams{
		AppName:                s.opts.HcpAppName,
		LocationOrganizationID: s.opts.HcpOrganizationID,
		LocationProjectID:      s.opts.HcpProjectID,
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

func (s *hashicorpClient) Delete(ctx context.Context, key string) error {
	params := &vault.DeleteAppSecretParams{
		AppName:                s.opts.HcpAppName,
		LocationOrganizationID: s.opts.HcpOrganizationID,
		LocationProjectID:      s.opts.HcpProjectID,
		Context:                ctx,
		SecretName:             key,
	}

	_, err := s.http.DeleteAppSecret(params, nil)
	return err
}
