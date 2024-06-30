package circleci

func (s *CircleciClient) Pipelines(org_slug string) (*CircleciPipelineResponse, error) {
	path := "/pipeline?org-slug=" + org_slug
	result := &CircleciPipelineResponse{}
	err := s.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *CircleciClient) PipelineWorkflows(pipelineId string) (*PipelineWorkflowResponse, error) {
	path := "pipeline/" + pipelineId + "/workflow"
	result := &PipelineWorkflowResponse{}
	err := s.http.Get(path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type CircleciPipelineResponse struct {
	Items         []CircleciPipeline `json:"items"`
	NextPageToken string             `json:"next_page_token"`
}

type CircleciPipeline struct {
	ID                string            `json:"id"`
	Errors            []CircleciError   `json:"errors"`
	ProjectSlug       string            `json:"project_slug"`
	UpdatedAt         string            `json:"updated_at"`
	Number            string            `json:"number"`
	TriggerParameters TriggerParameters `json:"trigger_parameters"`
	State             string            `json:"state"`
	CreatedAt         string            `json:"created_at"`
	Trigger           Trigger           `json:"trigger"`
	Vcs               Vcs               `json:"vcs"`
}

type CircleciError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Trigger struct {
	Type       string `json:"type"`
	ReceivedAt string `json:"received_at"`
	Actor      Actor  `json:"actor"`
}

type Actor struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}

type TriggerParameters struct {
	Property1 string `json:"property1"`
	Property2 string `json:"property2"`
}

type Vcs struct {
	ProviderName        string `json:"provider_name"`
	TargetRepositoryURL string `json:"target_repository_url"`
	Branch              string `json:"branch"`
	ReviewID            string `json:"review_id"`
	ReviewURL           string `json:"review_url"`
	Revision            string `json:"revision"`
	Tag                 string `json:"tag"`
	Commit              Commit `json:"commit"`
	OriginRepositoryURL string `json:"origin_repository_url"`
}

type Commit struct {
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type PipelineWorkflowResponse struct {
	Items         []PipelineWorkflowItem `json:"items"`
	NextPageToken string                 `json:"next_page_token"`
	Message       string                 `json:"message"`
}

type PipelineWorkflowItem struct {
	PipelineID     string `json:"pipeline_id"`
	CanceledBy     string `json:"canceled_by"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	ProjectSlug    string `json:"project_slug"`
	ErroredBy      string `json:"errored_by"`
	Tag            string `json:"tag"`
	Status         string `json:"status"`
	StartedBy      string `json:"started_by"`
	PipelineNumber string `json:"pipeline_number"`
	CreatedAt      string `json:"created_at"`
	StoppedAt      string `json:"stopped_at"`
}
