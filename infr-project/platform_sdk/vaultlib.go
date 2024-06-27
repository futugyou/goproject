package sdk

import (
	"context"
	"log"
	"os"

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

func (s *VaultClient) GetAppSecret(secretName string) (*vault.OpenAppSecretOK, error) {
	params := &vault.OpenAppSecretParams{
		AppName:                os.Getenv("HCP_APP_NAME"),
		LocationOrganizationID: os.Getenv("HCP_ORGANIZATION_ID"),
		LocationProjectID:      os.Getenv("HCP_PROJECT_ID"),
		SecretName:             secretName,
		Context:                context.Background(),
	}

	return s.http.OpenAppSecret(params, nil)
}
