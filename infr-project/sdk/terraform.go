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

func (s *TerraformClient) CheckWorkspace(name string) (*tfe.Workspace, error) {
	ctx := context.Background()

	w, err := s.client.Workspaces.Read(ctx, organization, name)
	if err != nil && err.Error() != "resource not found" {
		return nil, err
	}

	if w != nil {
		return w, nil
	}

	w, err = s.client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{
		Name: tfe.String(name),
	})
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (s *TerraformClient) CreateConfigurationVersions(workspaceID string, path string) (*tfe.ConfigurationVersion, error) {
	ctx := context.Background()

	cv, err := s.client.ConfigurationVersions.Create(ctx, workspaceID, tfe.ConfigurationVersionCreateOptions{
		AutoQueueRuns: tfe.Bool(false),
	})

	if err != nil {
		return nil, err
	}

	return cv, s.client.ConfigurationVersions.Upload(ctx, cv.UploadURL, path)
}
