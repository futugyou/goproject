package circleci

type WorkflowService service

func (s *WorkflowService) GetWorkflow(id string) (*JobDetailInfo, error) {
	path := "/workflow/" + id
	result := &JobDetailInfo{}
	err := s.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkflowService) ApproveJob(id string, approval_request_id string) (*BaseResponse, error) {
	path := "/workflow/" + id + "/approve/" + approval_request_id
	result := &BaseResponse{}
	err := s.client.http.Post(path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkflowService) CancelWorkflow(id string) (*BaseResponse, error) {
	path := "/workflow/" + id + "/cancel"
	result := &BaseResponse{}
	err := s.client.http.Post(path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkflowService) GetWorkflowJobs(id string) (*WorkflowJobsResponse, error) {
	path := "/workflow/" + id + "/job"
	result := &WorkflowJobsResponse{}
	err := s.client.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *WorkflowService) RerunWorkflow(id string) (*RerunWorkflowInfo, error) {
	path := "/workflow/" + id + "/rerun"
	result := &RerunWorkflowInfo{}
	err := s.client.http.Post(path, nil, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type RerunWorkflowInfo struct {
	WorkflowID string  `json:"workflow_id"`
	Message    *string `json:"message,omitempty"`
}

type WorkflowJobsResponse struct {
	Items         []WorkflowJob `json:"items"`
	NextPageToken string        `json:"next_page_token"`
	Message       *string       `json:"message,omitempty"`
}

type WorkflowJob struct {
	CanceledBy        string   `json:"canceled_by"`
	Dependencies      []string `json:"dependencies"`
	JobNumber         int64    `json:"job_number"`
	ID                string   `json:"id"`
	StartedAt         string   `json:"started_at"`
	Name              string   `json:"name"`
	ApprovedBy        string   `json:"approved_by"`
	ProjectSlug       string   `json:"project_slug"`
	Status            string   `json:"status"`
	Type              string   `json:"type"`
	StoppedAt         string   `json:"stopped_at"`
	ApprovalRequestID string   `json:"approval_request_id"`
}

type WorkflowInfo struct {
	PipelineID     string  `json:"pipeline_id"`
	CanceledBy     string  `json:"canceled_by"`
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	ProjectSlug    string  `json:"project_slug"`
	ErroredBy      string  `json:"errored_by"`
	Tag            string  `json:"tag"`
	Status         string  `json:"status"`
	StartedBy      string  `json:"started_by"`
	PipelineNumber string  `json:"pipeline_number"`
	CreatedAt      string  `json:"created_at"`
	StoppedAt      string  `json:"stopped_at"`
	Message        *string `json:"message,omitempty"`
}
