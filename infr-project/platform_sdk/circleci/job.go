package circleci

func (s *CircleciClient) GetJobDetails(project_slug string, job_number string) (*JobDetailInfo, error) {
	path := "/project/" + project_slug + "/job/" + job_number
	result := &JobDetailInfo{}
	err := s.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type JobDetailInfo struct {
	WebURL         string         `json:"web_url"`
	Project        Project        `json:"project"`
	ParallelRuns   []ParallelRun  `json:"parallel_runs"`
	StartedAt      string         `json:"started_at"`
	LatestWorkflow LatestWorkflow `json:"latest_workflow"`
	Name           string         `json:"name"`
	Executor       Executor       `json:"executor"`
	Parallelism    int64          `json:"parallelism"`
	Status         string         `json:"status"`
	Number         int64          `json:"number"`
	Pipeline       Pipeline       `json:"pipeline"`
	Duration       int64          `json:"duration"`
	CreatedAt      string         `json:"created_at"`
	Messages       []Message      `json:"messages"`
	Contexts       []Organization `json:"contexts"`
	Organization   Organization   `json:"organization"`
	QueuedAt       string         `json:"queued_at"`
	StoppedAt      string         `json:"stopped_at"`
	Message        *string        `json:"message"`
}

type Organization struct {
	Name string `json:"name"`
}

type Executor struct {
	ResourceClass string `json:"resource_class"`
	Type          string `json:"type"`
}

type LatestWorkflow struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
}

type ParallelRun struct {
	Index  int64  `json:"index"`
	Status string `json:"status"`
}

type Pipeline struct {
	ID string `json:"id"`
}

type Project struct {
	ID          string `json:"id"`
	Slug        string `json:"slug"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
}
