package sdk

import (
	"context"

	tfe "github.com/hashicorp/go-tfe"
)

const (
	tfcAPIBaseURL = "https://app.terraform.io/api/v2"
	organization  = "futugyousuzu"
	workspace     = "infr-project"
)

type TerraformClient struct {
	client *tfe.Client
}

func NewTerraformClient(token string) (*TerraformClient, error) {
	config := &tfe.Config{
		Token:             token,
		RetryServerErrors: true,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &TerraformClient{
		client: client,
	}, nil
}

func (s *TerraformClient) CheckWorkspace(name string) (string, error) {
	ctx := context.Background()

	w, err := s.client.Workspaces.Read(ctx, organization, name)
	if err != nil && err.Error() != "resource not found" {
		return "", err
	}

	if w != nil {
		return w.Name, nil
	}

	w, err = s.client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{
		Name: tfe.String(name),
	})
	if err != nil {
		return "", err
	}

	return w.Name, nil
}
