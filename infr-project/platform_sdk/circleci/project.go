package circleci

func (s *CircleciClient) CreateProject(org_slug string, name string) (*CreateProjectResponse, error) {
	path := "/project/" + org_slug + "/" + name

	result := &CreateProjectResponse{}
	if err := s.http.Post(path, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetProject(org_slug string, name string) (*ProjectInfo, error) {
	path := "/project/" + org_slug + "/" + name

	result := &ProjectInfo{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) CreateCheckoutKey(project_slug string, keyType string) (*CheckoutKey, error) {
	path := "/project/" + project_slug + "/checkout-key"

	request := &struct {
		Type string `json:"type"`
	}{
		Type: keyType,
	}
	result := &CheckoutKey{}
	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetCheckoutKey(project_slug string, digest string) (*CheckoutKeyList, error) {
	path := "/project/" + project_slug + "/checkout-key?digest=" + digest

	result := &CheckoutKeyList{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) DeleteCheckoutKey(project_slug string, fingerprint string) (*BaseResponse, error) {
	path := "/project/" + project_slug + "/checkout-key/" + fingerprint
	result := &BaseResponse{}
	if err := s.http.Delete(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetEnvironmentVariables(project_slug string) (*EnvironmentVariableList, error) {
	path := "/project/" + project_slug + "/envvar"

	result := &EnvironmentVariableList{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) CreateEnvironmentVariables(project_slug string, name string, value string) (*EnvironmentVariableInfo, error) {
	path := "/project/" + project_slug + "/envvar"
	request := struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}{
		Name:  name,
		Value: value,
	}
	result := &EnvironmentVariableInfo{}
	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) DeleteEnvironmentVariables(project_slug string, name string) (*BaseResponse, error) {
	path := "/project/" + project_slug + "/envvar/" + name

	result := &BaseResponse{}
	if err := s.http.Delete(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetMaskedEnvironmentVariable(project_slug string, name string) (*EnvironmentVariableInfo, error) {
	path := "/project/" + project_slug + "/envvar/" + name

	result := &EnvironmentVariableInfo{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
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
