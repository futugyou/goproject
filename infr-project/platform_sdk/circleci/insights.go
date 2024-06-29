package circleci

// Example: workflow-names=A single workflow name: ?workflow-names=build-test-deploy
// or for multiple workflow names: ?workflow-names=build&workflow-names=test-and-deploy.
func (s *CircleciClient) GetMetricsTrends(project_slug string, workflow_names []string) (*MetricsTrends, error) {
	path := "/insights/pages/" + project_slug + "/summary"
	for i := 0; i < len(workflow_names); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("workflow-names=" + workflow_names[i])
	}

	result := &MetricsTrends{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetJobTimeseriesData(project_slug string, workflow_names string) (*JobTimeseriesDataResponse, error) {
	path := "/insights/time-series/" + project_slug + "/workflows/jobs"

	result := &JobTimeseriesDataResponse{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetOrgMetricsTrends(org_slug string, project_names []string) (*OrgMetricsTrends, error) {
	path := "/insights/pages/" + org_slug + "/summary"
	for i := 0; i < len(project_names); i++ {
		if i == 0 {
			path += "?"
		} else {
			path += "&"
		}
		path += ("project-names=" + project_names[i])
	}

	result := &OrgMetricsTrends{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetBranches(project_slug string, workflow_ame string) (*ProjectBranches, error) {
	path := "/insights/" + project_slug + "/branches"
	if len(workflow_ame) > 0 {
		path += ("?workflow-name=" + workflow_ame)
	}

	result := &ProjectBranches{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) GetFlakyTests(project_slug string) (*FlakyTestsResponse, error) {
	path := "/insights/" + project_slug + "/flaky-tests"
	result := &FlakyTestsResponse{}
	if err := s.http.Get(path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type FlakyTestsResponse struct {
	FlakyTests      []FlakyTest `json:"flaky-tests"`
	TotalFlakyTests int64       `json:"total-flaky-tests"`
	Message         *string     `json:"message,omitempty"`
}

type FlakyTest struct {
	TimeWasted        int64       `json:"time-wasted"`
	WorkflowCreatedAt string      `json:"workflow-created-at"`
	WorkflowID        interface{} `json:"workflow-id"`
	Classname         string      `json:"classname"`
	PipelineNumber    int64       `json:"pipeline-number"`
	WorkflowName      string      `json:"workflow-name"`
	TestName          string      `json:"test-name"`
	JobName           string      `json:"job-name"`
	JobNumber         int64       `json:"job-number"`
	TimesFlaked       int64       `json:"times-flaked"`
	Source            string      `json:"source"`
	File              string      `json:"file"`
}

type ProjectBranches struct {
	OrgID     string   `json:"org_id"`
	ProjectID string   `json:"project_id"`
	Branches  []string `json:"branches"`
	Message   *string  `json:"message,omitempty"`
}

type OrgMetricsTrends struct {
	OrgData        OrgData           `json:"org_data"`
	OrgProjectData []OrgProjectDatum `json:"org_project_data"`
	AllProjects    []string          `json:"all_projects"`
	Message        *string           `json:"message,omitempty"`
}

type OrgData struct {
	Metrics Metrics `json:"metrics"`
	Trends  Metrics `json:"trends"`
}

type OrgProjectDatum struct {
	ProjectName string  `json:"project_name"`
	Metrics     Metrics `json:"metrics"`
	Trends      Metrics `json:"trends"`
}

type JobTimeseriesDataResponse struct {
	NextPageToken string              `json:"next_page_token"`
	Items         []JobTimeseriesData `json:"items"`
	Message       *string             `json:"message,omitempty"`
}

type JobTimeseriesData struct {
	Name         string  `json:"name"`
	MinStartedAt string  `json:"min_started_at"`
	MaxEndedAt   string  `json:"max_ended_at"`
	Timestamp    string  `json:"timestamp"`
	Metrics      Metrics `json:"metrics"`
}

type Metrics struct {
	TotalRuns         int64           `json:"total_runs"`
	FailedRuns        int64           `json:"failed_runs"`
	SuccessfulRuns    int64           `json:"successful_runs"`
	Throughput        int64           `json:"throughput"`
	MedianCreditsUsed int64           `json:"median_credits_used"`
	TotalCreditsUsed  int64           `json:"total_credits_used"`
	DurationMetrics   DurationMetrics `json:"duration_metrics"`
}

type DurationMetrics struct {
	Min    int64 `json:"min"`
	Median int64 `json:"median"`
	Max    int64 `json:"max"`
	P95    int64 `json:"p95"`
	Total  int64 `json:"total"`
}

type MetricsTrends struct {
	OrgID                     string                 `json:"org_id"`
	ProjectID                 string                 `json:"project_id"`
	ProjectData               ProjectData            `json:"project_data"`
	ProjectWorkflowData       []ProjectWorkflowDatum `json:"project_workflow_data"`
	ProjectWorkflowBranchData []ProjectWorkflowDatum `json:"project_workflow_branch_data"`
	AllBranches               []string               `json:"all_branches"`
	AllWorkflows              []string               `json:"all_workflows"`
	Message                   *string                `json:"message,omitempty"`
}

type ProjectData struct {
	Metrics ProjectDataMetrics `json:"metrics"`
	Trends  ProjectDataMetrics `json:"trends"`
}

type ProjectDataMetrics struct {
	TotalRuns         int64 `json:"total_runs"`
	TotalDurationSecs int64 `json:"total_duration_secs"`
	TotalCreditsUsed  int64 `json:"total_credits_used"`
	SuccessRate       int64 `json:"success_rate"`
	Throughput        int64 `json:"throughput"`
}

type ProjectWorkflowDatum struct {
	WorkflowName string                            `json:"workflow_name"`
	Branch       *string                           `json:"branch,omitempty"`
	Metrics      ProjectWorkflowBranchDatumMetrics `json:"metrics"`
	Trends       ProjectWorkflowBranchDatumMetrics `json:"trends"`
}

type ProjectWorkflowBranchDatumMetrics struct {
	TotalCreditsUsed int64 `json:"total_credits_used"`
	P95DurationSecs  int64 `json:"p95_duration_secs"`
	TotalRuns        int64 `json:"total_runs"`
	SuccessRate      int64 `json:"success_rate"`
}
