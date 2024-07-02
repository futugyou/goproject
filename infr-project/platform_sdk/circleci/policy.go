package circleci

func (s *CircleciClient) DecisionAuditLogs(ownerID string, context string) ([]AuditLogInfo, error) {
	path := "/owner/" + ownerID + "/context/" + context + "/decision"
	result := []AuditLogInfo{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}

	return result, nil
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
