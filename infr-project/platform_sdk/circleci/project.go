package circleci

import "log"

func (s *CircleciClient) CreateProject(org_slug string, name string) CreateProjectResponse {
	path := "/project/" + org_slug + "/" + name

	result := CreateProjectResponse{}
	err := s.http.Post(path, nil, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) GetProject(org_slug string, name string) ProjectInfo {
	path := "/project/" + org_slug + "/" + name

	result := ProjectInfo{}
	err := s.http.Get(path, &result)

	if err != nil {
		log.Println(err.Error())
		return result
	}
	return result
}

func (s *CircleciClient) CreateCheckoutKey(project_slug string, keyType string) (*CheckoutKey, error) {
	path := "/project/" + project_slug + "/checkout-key"

	request := &struct {
		Type string `json:"type"`
	}{
		Type: keyType,
	}
	result := &CheckoutKey{}
	err := s.http.Post(path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
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

type VcsInfo struct {
	VcsURL        string `json:"vcs_url"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"default_branch"`
}
