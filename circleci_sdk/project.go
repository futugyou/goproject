package circleci

type ProjectService service

// org_slug include provider and organization
// project_slug include provider and organization and project
// eg. org_slug gh/CircleCI-Public
// eg. project_slug gh/CircleCI-Public/api-preview-docs
func (s *ProjectService) CreateProject(org_slug string, name string) (*CreateProjectResponse, error) {
	path := "/project/" + org_slug + "/" + name

	result := &CreateProjectResponse{}
	if err := s.client.http.Post(path, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetProject(org_slug string, name string) (*ProjectInfo, error) {
	path := "/project/" + org_slug + "/" + name

	result := &ProjectInfo{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) CreateCheckoutKey(project_slug string, keyType string) (*CheckoutKey, error) {
	path := "/project/" + project_slug + "/checkout-key"

	request := &struct {
		Type string `json:"type"`
	}{
		Type: keyType,
	}
	result := &CheckoutKey{}
	if err := s.client.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetCheckoutKey(project_slug string, digest string) (*CheckoutKeyList, error) {
	path := "/project/" + project_slug + "/checkout-key?digest=" + digest

	result := &CheckoutKeyList{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) DeleteCheckoutKey(project_slug string, fingerprint string) (*BaseResponse, error) {
	path := "/project/" + project_slug + "/checkout-key/" + fingerprint
	result := &BaseResponse{}
	if err := s.client.http.Delete(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetEnvironmentVariables(project_slug string) (*EnvironmentVariableList, error) {
	path := "/project/" + project_slug + "/envvar"

	result := &EnvironmentVariableList{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) CreateEnvironmentVariables(project_slug string, name string, value string) (*EnvironmentVariableInfo, error) {
	path := "/project/" + project_slug + "/envvar"
	request := struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  name,
		Value: value,
	}
	result := &EnvironmentVariableInfo{}
	if err := s.client.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) DeleteEnvironmentVariables(project_slug string, name string) (*BaseResponse, error) {
	path := "/project/" + project_slug + "/envvar/" + name

	result := &BaseResponse{}
	if err := s.client.http.Delete(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetMaskedEnvironmentVariable(project_slug string, name string) (*EnvironmentVariableInfo, error) {
	path := "/project/" + project_slug + "/envvar/" + name

	result := &EnvironmentVariableInfo{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) GetProjectSettings(project_slug string) (*ProjectSettingList, error) {
	path := "/project/" + project_slug + "/settings"

	result := &ProjectSettingList{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ProjectService) UpdateProjectSettings(project_slug string, advanced Advanced) (*ProjectSettingList, error) {
	path := "/project/" + project_slug + "/settings"

	result := &ProjectSettingList{}
	if err := s.client.http.Patch(path, advanced, result); err != nil {
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

type VcsInfo struct {
	VcsURL        string `json:"vcs_url"`
	Provider      string `json:"provider"`
	DefaultBranch string `json:"default_branch"`
}
