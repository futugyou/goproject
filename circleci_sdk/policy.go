package circleci

import (
	"context"
	"fmt"
	"net/url"
)

type PolicyService service

func (s *PolicyService) DecisionAuditLogs(ctx context.Context, ownerID string, context string) ([]AuditLogInfo, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision", ownerID, context)
	result := []AuditLogInfo{}
	if err := s.client.http.Get(ctx, path, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) MakeDecision(ctx context.Context, ownerID string, context string, input string, metadata interface{}) (*MakeDecisionResponse, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision", ownerID, context)
	request := &struct {
		Input    string      `json:"input"`
		Metadata interface{} `json:"metadata"`
	}{
		Input:    input,
		Metadata: metadata,
	}
	result := &MakeDecisionResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) GetDecisionSettings(ctx context.Context, ownerID string, context string) (*DecisionSetting, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision/settings", ownerID, context)
	result := &DecisionSetting{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) SetDecisionSettings(ctx context.Context, ownerID string, context string, enabled bool) (*DecisionSetting, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision/settings", ownerID, context)
	request := &struct {
		Enabled bool `json:"enabled"`
	}{
		Enabled: enabled,
	}
	result := &DecisionSetting{}
	if err := s.client.http.Patch(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) DecisionAuditLogById(ctx context.Context, ownerID string, context string, decisionID string) (*AuditLogInfo, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision/%s", ownerID, context, decisionID)
	result := &AuditLogInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) PolicyBundleLogById(ctx context.Context, ownerID string, context string, decisionID string) (map[string]PolicyProperty, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/decision/%s/policy-bundle", ownerID, context, decisionID)
	result := make(map[string]PolicyProperty)
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) PolicyBundles(ctx context.Context, ownerID string, context string) (map[string]PolicyProperty, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/policy-bundle", ownerID, context)
	result := make(map[string]PolicyProperty)
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) CreatesPolicyBundle(ctx context.Context, ownerID string, context string, dry bool, policies map[string]string) (*CreatesPolicyBundleResponse, error) {
	u := &url.URL{
		Path: fmt.Sprintf("/owner/%s/context/%s/policy-bundle", ownerID, context),
	}
	params := url.Values{}
	params.Add("dry", fmt.Sprintf("%t", dry))
	u.RawQuery = params.Encode()
	path := u.String()

	request := &struct {
		Policies map[string]string `json:"policies"`
	}{
		Policies: policies,
	}
	result := &CreatesPolicyBundleResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PolicyService) RetrievesPolicyDocument(ctx context.Context, ownerID string, context string, policyName string) (*PolicyProperty, error) {
	path := fmt.Sprintf("/owner/%s/context/%s/policy-bundle/%s", ownerID, context, policyName)
	result := &PolicyProperty{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
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
