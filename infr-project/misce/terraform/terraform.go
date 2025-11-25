package terraform

import (
	"bytes"
	"context"

	"github.com/hashicorp/go-slug"
	tfe "github.com/hashicorp/go-tfe"
)

type TerraformClient struct {
	client        *tfe.Client
	tfcAPIBaseURL string
	organization  string
	workspace     string
}

func NewTerraformClient(token string, tfcAPIBaseURL string, organization string, workspace string) (*TerraformClient, error) {
	config := &tfe.Config{
		Token:             token,
		RetryServerErrors: true,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &TerraformClient{
		client:        client,
		tfcAPIBaseURL: tfcAPIBaseURL,
		organization:  organization,
		workspace:     workspace,
	}, nil
}

func (s *TerraformClient) CheckWorkspace(name string) (*tfe.Workspace, error) {
	ctx := context.Background()

	w, err := s.client.Workspaces.Read(ctx, s.organization, name)
	if err != nil && err.Error() != "resource not found" {
		return nil, err
	}

	if w != nil {
		return w, nil
	}

	w, err = s.client.Workspaces.Create(ctx, s.organization, tfe.WorkspaceCreateOptions{
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

	packer, err := slug.NewPacker(
		slug.DereferenceSymlinks(),    // dereferences symlinks
		slug.ApplyTerraformIgnore(),   // ignores paths specified in .terraformignore
		slug.AllowSymlinkTarget(path), // allow certain symlink target paths
	)

	if err != nil {
		return nil, err
	}

	rawConfig := bytes.NewBuffer(nil)
	_, err = packer.Pack(path, rawConfig)
	if err != nil {
		return nil, err
	}

	return cv, s.client.ConfigurationVersions.UploadTarGzip(ctx, cv.UploadURL, rawConfig)
}

func (s *TerraformClient) CreateRun(workspace *tfe.Workspace, planOnly bool) (*tfe.Run, error) {
	ctx := context.Background()
	options := tfe.RunCreateOptions{
		Workspace: workspace,
		PlanOnly:  tfe.Bool(planOnly),
	}
	return s.client.Runs.Create(ctx, options)
}

func (s *TerraformClient) ApplyRun(runID string) error {
	ctx := context.Background()
	options := tfe.RunApplyOptions{}
	return s.client.Runs.Apply(ctx, runID, options)
}
