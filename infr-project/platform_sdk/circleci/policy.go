package circleci

func (s *CircleciClient) DecisionAuditLogs(ownerID string, context string) ([]AuditLogInfo, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision"
	result := []AuditLogInfo{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) MakeDecision(ownerID string, context string, input string, metadata interface{}) (*MakeDecisionResponse, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision"
	request := &struct {
		Input    string      `json:"input"`
		Metadata interface{} `json:"metadata"`
	}{
		Input:    input,
		Metadata: metadata,
	}
	result := &MakeDecisionResponse{}
	if err := s.http.Post(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) GetDecisionSettings(ownerID string, context string) (*DecisionSetting, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/settings"
	result := &DecisionSetting{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CircleciClient) SetDecisionSettings(ownerID string, context string, enabled bool) (*DecisionSetting, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/settings"
	request := &struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: enabled,
	}
	result := &DecisionSetting{}
	if err := s.http.Patch(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

type DecisionSetting struct {
	Enabled bool    `json:"enabled"`
	Error   *string `json:"error"`
}

type MakeDecisionResponse struct {
	EnabledRules []string  `json:"enabled_rules"`
	HardFailures []Failure `json:"hard_failures"`
	Reason       string    `json:"reason"`
	SoftFailures []Failure `json:"soft_failures"`
	Status       string    `json:"status"`
	Error        *string   `json:"error"`
}

type AuditLogInfo struct {
	CreatedAt   string   `json:"created_at"`
	Decision    Decision `json:"decision"`
	ID          string   `json:"id"`
	Metadata    Metadata `json:"metadata"`
	Policies    Policies `json:"policies"`
	TimeTakenMS int64    `json:"time_taken_ms"`
}

type Decision struct {
	EnabledRules []string  `json:"enabled_rules"`
	HardFailures []Failure `json:"hard_failures"`
	Reason       string    `json:"reason"`
	SoftFailures []Failure `json:"soft_failures"`
	Status       string    `json:"status"`
}

type Failure struct {
	Reason string `json:"reason"`
	Rule   string `json:"rule"`
}

type Metadata struct {
	BuildNumber int64    `json:"build_number"`
	ProjectID   string   `json:"project_id"`
	SSHRerun    bool     `json:"ssh_rerun"`
	Vcs         AuditVcs `json:"vcs"`
}

type AuditVcs struct {
	Branch              string `json:"branch"`
	OriginRepositoryURL string `json:"origin_repository_url"`
	ReleaseTag          string `json:"release_tag"`
	TargetRepositoryURL string `json:"target_repository_url"`
}

type Policies struct {
	PolicyName1 string `json:"policy_name1"`
	PolicyName2 string `json:"policy_name2"`
}
