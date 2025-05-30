package circleci

import (
	"context"
	"fmt"
	"net/url"
)

type ProjectService service

// org_slug include provider and organization
// project_slug include provider and organization and project
// eg. org_slug gh/CircleCI-Public
// eg. project_slug gh/CircleCI-Public/api-preview-docs
func (s *ProjectService) CreateProject(ctx context.Context, org_slug string, name string) (*CreateProjectResponse, error) {
	path := fmt.Sprintf("/project/%s/%s", org_slug, name)

	result := &CreateProjectResponse{}
	if err := s.client.http.Post(ctx, path, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetProject(ctx context.Context, org_slug string, name string) (*ProjectInfo, error) {
	path := fmt.Sprintf("/project/%s/%s", org_slug, name)

	result := &ProjectInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) ListProject(ctx context.Context) ([]ProjectListItem, error) {
	path := "/projects"

	result := []ProjectListItem{}
	if err := s.client.http.Get(ctx, path, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) CreateCheckoutKey(ctx context.Context, project_slug string, keyType string) (*CheckoutKey, error) {
	path := fmt.Sprintf("/project/%s/checkout-key", project_slug)

	request := &struct {
		Type string `json:"type"`
	}{
		Type: keyType,
	}
	result := &CheckoutKey{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetCheckoutKey(ctx context.Context, project_slug string, digest string) (*CheckoutKeyList, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/project/%s/checkout-key", project_slug),
	}
	params := url.Values{}
	params.Add("digest", digest)
	u.RawQuery = params.Encode()
	path := u.String()
	result := &CheckoutKeyList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) DeleteCheckoutKey(ctx context.Context, project_slug string, fingerprint string) (*BaseResponse, error) {
	path := fmt.Sprintf("/project/%s/checkout-key/%s", project_slug, fingerprint)

	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetEnvironmentVariables(ctx context.Context, project_slug string) (*EnvironmentVariableList, error) {
	path := fmt.Sprintf("/project/%s/envvar", project_slug)

	result := &EnvironmentVariableList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) CreateEnvironmentVariables(ctx context.Context, project_slug string, name string, value string) (*EnvironmentVariableInfo, error) {
	path := fmt.Sprintf("/project/%s/envvar", project_slug)
	request := struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  name,
		Value: value,
	}
	result := &EnvironmentVariableInfo{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) DeleteEnvironmentVariables(ctx context.Context, project_slug string, name string) (*BaseResponse, error) {
	path := fmt.Sprintf("/project/%s/envvar/%s", project_slug, name)

	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetMaskedEnvironmentVariable(ctx context.Context, project_slug string, name string) (*EnvironmentVariableInfo, error) {
	path := fmt.Sprintf("/project/%s/envvar/%s", project_slug, name)

	result := &EnvironmentVariableInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetProjectSettings(ctx context.Context, project_slug string) (*ProjectSettingList, error) {
	path := fmt.Sprintf("/project/%s/settings", project_slug)

	result := &ProjectSettingList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) UpdateProjectSettings(ctx context.Context, project_slug string, advanced Advanced) (*ProjectSettingList, error) {
	path := fmt.Sprintf("/project/%s/settings", project_slug)

	result := &ProjectSettingList{}
	if err := s.client.http.Patch(ctx, path, advanced, result); err != nil {
		return nil, err
	}
	return result, nil
}

type ProjectSettingList struct {
	Advanced Advanced `json:"advanced"`
	Message  *string  `json:"message,omitempty"`
}

type EnvironmentVariableList struct {
	Items         []EnvironmentVariableInfo `json:"items"`
	NextPageToken string                    `json:"next_page_token"`
	Message       *string                   `json:"message,omitempty"`
}

type EnvironmentVariableInfo struct {
	Name      string  `json:"name"`
	Value     string  `json:"value"`
	CreatedAt string  `json:"created-at"`
	Message   *string `json:"message,omitempty"`
}

type CheckoutKeyList struct {
	Items         []CheckoutKey `json:"items"`
	NextPageToken string        `json:"next_page_token"`
	Message       *string       `json:"message,omitempty"`
}

type CheckoutKey struct {
	PublicKey   string `json:"public-key"`
	Type        string `json:"type"`
	Fingerprint string `json:"fingerprint"`
	Preferred   bool   `json:"preferred"`
	CreatedAt   string `json:"created-at"`
}

type CreateProjectResponse struct {
	Advanced Advanced `json:"advanced"`
	Message  *string  `json:"message,omitempty"`
}

type Advanced struct {
	AutocancelBuilds           bool     `json:"autocancel_builds"`
	BuildForkPrs               bool     `json:"build_fork_prs"`
	BuildPrsOnly               bool     `json:"build_prs_only"`
	DisableSSH                 bool     `json:"disable_ssh"`
	ForksReceiveSecretEnvVars  bool     `json:"forks_receive_secret_env_vars"`
	OSS                        bool     `json:"oss"`
	SetGithubStatus            bool     `json:"set_github_status"`
	SetupWorkflows             bool     `json:"setup_workflows"`
	WriteSettingsRequiresAdmin bool     `json:"write_settings_requires_admin"`
	PROnlyBranchOverrides      []string `json:"pr_only_branch_overrides"`
}

type ProjectInfo struct {
	Slug             string  `json:"slug"`
	Name             string  `json:"name"`
	ID               string  `json:"id"`
	OrganizationName string  `json:"organization_name"`
	OrganizationSlug string  `json:"organization_slug"`
	OrganizationID   string  `json:"organization_id"`
	VcsInfo          VcsInfo `json:"vcs_info"`
	Message          *string `json:"message,omitempty"`
}

func (r *ProjectInfo) GetMessage() string {
	if r == nil || r.Message == nil {
		return ""
	}
	return *r.Message
}

type VcsInfo struct {
	VcsURL        string `json:"vcs_url"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"default_branch"`
}

type ProjectListItem struct {
	OSS           bool                  `json:"oss"`
	Reponame      string                `json:"reponame"`
	Username      string                `json:"username"`
	HasUsableKey  bool                  `json:"has_usable_key"`
	VcsType       string                `json:"vcs_type"`
	Language      interface{}           `json:"language"`
	VcsURL        string                `json:"vcs_url"`
	Following     bool                  `json:"following"`
	DefaultBranch string                `json:"default_branch"`
	Branches      map[string]BranchInfo `json:"branches"`
}

type BranchInfo struct {
	RunningBuilds            []BranchBuildInfo             `json:"running_builds"`
	RecentBuilds             []BranchBuildInfo             `json:"recent_builds"`
	IsUsingWorkflows         bool                          `json:"is_using_workflows"`
	PusherLogins             []string                      `json:"pusher_logins"`
	LastSuccess              BranchBuildInfo               `json:"last_success"`
	LastNonSuccess           BranchBuildInfo               `json:"last_non_success"`
	LatestWorkflows          map[string]BranchWorkflowInfo `json:"latest_workflows"`
	LatestCompletedWorkflows map[string]BranchWorkflowInfo `json:"latest_completed_workflows"`
}

type BranchBuildInfo struct {
	Status        string `json:"status"`
	Outcome       string `json:"outcome"`
	BuildNum      int64  `json:"build_num"`
	VcsRevision   string `json:"vcs_revision"`
	PushedAt      string `json:"pushed_at"`
	AddedAt       string `json:"added_at"`
	IsWorkflowJob bool   `json:"is_workflow_job"`
	Is2_0_Job     bool   `json:"is_2_0_job"`
}

type BranchWorkflowInfo struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}
