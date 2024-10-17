package circleci

import "fmt"

type PolicyService service

func (s *PolicyService) DecisionAuditLogs(ownerID string, context string) ([]AuditLogInfo, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision"
	result := []AuditLogInfo{}
	if err := s.client.http.Get(path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) MakeDecision(ownerID string, context string, input string, metadata interface{}) (*MakeDecisionResponse, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision"
	request := &struct {
		Input    string      `json:"input"`
		Metadata interface{} `json:"metadata"`
	}{
		Input:    input,
		Metadata: metadata,
	}
	result := &MakeDecisionResponse{}
	if err := s.client.http.Post(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) GetDecisionSettings(ownerID string, context string) (*DecisionSetting, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/settings"
	result := &DecisionSetting{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) SetDecisionSettings(ownerID string, context string, enabled bool) (*DecisionSetting, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/settings"
	request := &struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: enabled,
	}
	result := &DecisionSetting{}
	if err := s.client.http.Patch(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) DecisionAuditLogById(ownerID string, context string, decisionID string) (*AuditLogInfo, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/" + decisionID
	result := &AuditLogInfo{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) PolicyBundleLogById(ownerID string, context string, decisionID string) (map[string]PolicyProperty, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision/" + decisionID + "/policy-bundle"
	result := make(map[string]PolicyProperty)
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) PolicyBundles(ownerID string, context string) (map[string]PolicyProperty, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/policy-bundle"
	result := make(map[string]PolicyProperty)
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) CreatesPolicyBundle(ownerID string, context string, dry bool, policies map[string]string) (*CreatesPolicyBundleResponse, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/policy-bundle?dry=" + fmt.Sprintf("%t", dry)
	request := &struct {
		Policies map[string]string `json:"policies"`
	}{
		Policies: policies,
	}
	result := &CreatesPolicyBundleResponse{}
	if err := s.client.http.Post(path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) RetrievesPolicyDocument(ownerID string, context string, policyName string) (*PolicyProperty, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/policy-bundle/" + policyName
	result := &PolicyProperty{}
	if err := s.client.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
}

type CreatesPolicyBundleResponse struct {
	Created  []string `json:"created"`
	Deleted  []string `json:"deleted"`
	Modified []string `json:"modified"`
	Error    *string  `json:"error"`
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
	CreatedAt   string            `json:"created_at"`
	Decision    Decision          `json:"decision"`
	ID          string            `json:"id"`
	Metadata    Metadata          `json:"metadata"`
	Policies    map[string]string `json:"policies"`
	TimeTakenMS int64             `json:"time_taken_ms"`
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

type PolicyProperty struct {
	Content   string  `json:"content"`
	CreatedAt string  `json:"created_at"`
	CreatedBy string  `json:"created_by"`
	Name      string  `json:"name"`
	Error     *string `json:"error"`
}
