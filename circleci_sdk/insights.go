package circleci

import "context"

type InsightsService service

// Example: workflow-names=A single workflow name: ?workflow-names=build-test-deploy
// or for multiple workflow names: ?workflow-names=build&workflow-names=test-and-deploy.
func (s *InsightsService) GetMetricsTrends(ctx context.Context, project_slug string, workflow_names []string) (*MetricsTrends, error) {
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
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetJobTimeseriesData(ctx context.Context, project_slug string, workflow_names string) (*JobTimeseriesDataResponse, error) {
	path := "/insights/time-series/" + project_slug + "/workflows/jobs"

	result := &JobTimeseriesDataResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetOrgMetricsTrends(ctx context.Context, org_slug string, project_names []string) (*OrgMetricsTrends, error) {
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
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetBranches(ctx context.Context, project_slug string, workflow_ame string) (*ProjectBranches, error) {
	path := "/insights/" + project_slug + "/branches"
	if len(workflow_ame) > 0 {
		path += ("?workflow-name=" + workflow_ame)
	}

	result := &ProjectBranches{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetFlakyTests(ctx context.Context, project_slug string) (*FlakyTestsResponse, error) {
	path := "/insights/" + project_slug + "/flaky-tests"
	result := &FlakyTestsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetMetricsForProjectWorkflows(ctx context.Context, project_slug string) (*WorkflowMetricsResponse, error) {
	path := "/insights/" + project_slug + "/workflows"
	result := &WorkflowMetricsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetRecentRuns(ctx context.Context, project_slug string, workflow_name string) (*RecentRunsResponse, error) {
	path := "/insights/" + project_slug + "/workflows/" + workflow_name
	result := &RecentRunsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetMetricsForProjectWorkflowJobs(ctx context.Context, project_slug string, workflow_name string) (*MetricsProjectWorkflowJobsResponse, error) {
	path := "/insights/" + project_slug + "/workflows/" + workflow_name + "/jobs"
	result := &MetricsProjectWorkflowJobsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetMetricsTrendsForWorkflows(ctx context.Context, project_slug string, workflow_name string) (*MetricsTrendsForWorkflowsResponse, error) {
	path := "/insights/" + project_slug + "/workflows/" + workflow_name + "/summary"

	result := &MetricsTrendsForWorkflowsResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *InsightsService) GetTestMetricsForProjectWorkflow(ctx context.Context, project_slug string, workflow_name string) (*TestMetrics, error) {
	path := "/insights/" + project_slug + "/workflows/" + workflow_name + "/test-metrics"
	result := &TestMetrics{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type TestMetrics struct {
	AverageTestCount     int64     `json:"average_test_count"`
	MostFailedTests      []Test    `json:"most_failed_tests"`
	MostFailedTestsExtra int64     `json:"most_failed_tests_extra"`
	SlowestTests         []Test    `json:"slowest_tests"`
	SlowestTestsExtra    int64     `json:"slowest_tests_extra"`
	TotalTestRuns        int64     `json:"total_test_runs"`
	TestRuns             []TestRun `json:"test_runs"`
	Message              *string   `json:"message,omitempty"`
}

type Test struct {
	P95Duration int64  `json:"p95_duration"`
	TotalRuns   int64  `json:"total_runs"`
	Classname   string `json:"classname"`
	FailedRuns  int64  `json:"failed_runs"`
	Flaky       bool   `json:"flaky"`
	Source      string `json:"source"`
	File        string `json:"file"`
	JobName     string `json:"job_name"`
	TestName    string `json:"test_name"`
}

type TestRun struct {
	PipelineNumber int64       `json:"pipeline_number"`
	WorkflowID     interface{} `json:"workflow_id"`
	SuccessRate    int64       `json:"success_rate"`
	TestCounts     TestCounts  `json:"test_counts"`
}

type TestCounts struct {
	Error   int64 `json:"error"`
	Failure int64 `json:"failure"`
	Skipped int64 `json:"skipped"`
	Success int64 `json:"success"`
	Total   int64 `json:"total"`
}

type MetricsTrendsForWorkflowsResponse struct {
	Metrics       Metrics  `json:"metrics"`
	Trends        Trends   `json:"trends"`
	WorkflowNames []string `json:"workflow_names"`
	Message       *string  `json:"message,omitempty"`
}

type Trends struct {
	TotalRuns          int64 `json:"total_runs"`
	FailedRuns         int64 `json:"failed_runs"`
	SuccessRate        int64 `json:"success_rate"`
	P95DurationSecs    int64 `json:"p95_duration_secs"`
	MedianDurationSecs int64 `json:"median_duration_secs"`
	TotalCreditsUsed   int64 `json:"total_credits_used"`
	Mttr               int64 `json:"mttr"`
	Throughput         int64 `json:"throughput"`
}

type MetricsProjectWorkflowJobsResponse struct {
	Items         []WorkflowMetrics `json:"items"`
	NextPageToken string            `json:"next_page_token"`
}

type RecentRunsResponse struct {
	Items         []RunsInfo `json:"items"`
	NextPageToken string     `json:"next_page_token"`
	Message       *string    `json:"message,omitempty"`
}

type RunsInfo struct {
	ID          string `json:"id"`
	Branch      string `json:"branch"`
	Duration    int64  `json:"duration"`
	CreatedAt   string `json:"created_at"`
	StoppedAt   string `json:"stopped_at"`
	CreditsUsed int64  `json:"credits_used"`
	Status      string `json:"status"`
	IsApproval  bool   `json:"is_approval"`
}

type WorkflowMetricsResponse struct {
	Items         []WorkflowMetrics `json:"items"`
	NextPageToken string            `json:"next_page_token"`
	Message       *string           `json:"message,omitempty"`
}

type WorkflowMetrics struct {
	Name        string  `json:"name"`
	Metrics     Metrics `json:"metrics"`
	WindowStart string  `json:"window_start"`
	WindowEnd   string  `json:"window_end"`
	ProjectID   string  `json:"project_id"`
}

type FlakyTestsResponse struct {
	FlakyTests      []FlakyTest `json:"flaky-tests"`
	TotalFlakyTests int64       `json:"total-flaky-tests"`
	Message         *string     `json:"message,omitempty"`
}

type FlakyTest struct {
	TimeWasted        int64  `json:"time-wasted"`
	WorkflowCreatedAt string `json:"workflow-created-at"`
	WorkflowID        string `json:"workflow-id"`
	Classname         string `json:"classname"`
	PipelineNumber    int64  `json:"pipeline-number"`
	WorkflowName      string `json:"workflow-name"`
	TestName          string `json:"test-name"`
	JobName           string `json:"job-name"`
	JobNumber         int64  `json:"job-number"`
	TimesFlaked       int64  `json:"times-flaked"`
	Source            string `json:"source"`
	File              string `json:"file"`
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
	Mttr              int64           `json:"mttr"`
	SuccessRate       int64           `json:"success_rate"`
	TotalRecoveries   int64           `json:"total_recoveries"`
}

type DurationMetrics struct {
	Min               int64 `json:"min"`
	Median            int64 `json:"median"`
	Max               int64 `json:"max"`
	P95               int64 `json:"p95"`
	Total             int64 `json:"total"`
	StandardDeviation int64 `json:"standard_deviation"`
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
