package circleci

import (
	"context"
	"fmt"
)

type PipelineService service

func (s *PipelineService) Pipelines(ctx context.Context, org_slug string) (*CircleciPipelineResponse, error) {
	path := fmt.Sprintf("/pipeline?org-slug=%s", org_slug)
	result := &CircleciPipelineResponse{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) GetPipelineById(ctx context.Context, pipelineid string) (*CircleciPipeline, error) {
	path := fmt.Sprintf("/pipeline/%s/config", pipelineid)
	result := &CircleciPipeline{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) GetPipelineConfiguration(ctx context.Context, pipelineid string) (*PipelineConfig, error) {
	path := fmt.Sprintf("/pipeline/%s", pipelineid)
	result := &PipelineConfig{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) PipelineWorkflows(ctx context.Context, pipelineId string) (*PipelineWorkflowResponse, error) {
	path := fmt.Sprintf("/pipeline/%s/workflow", pipelineId)
	result := &PipelineWorkflowResponse{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) ContinuePipeline(ctx context.Context, continuationKey string, configuration string, parameters interface{}) (*BaseResponse, error) {
	path := "/pipeline/continue"
	request := ContinuePipelineRequest{
		ContinuationKey: continuationKey,
		Configuration:   configuration,
		Parameters:      parameters,
	}
	result := &BaseResponse{}
	err := s.client.http.Post(ctx, path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) GetPipelinesByProject(ctx context.Context, project_slug string) (*CircleciPipelineResponse, error) {
	path := fmt.Sprintf("/project/%s/pipeline", project_slug)
	result := &CircleciPipelineResponse{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) TriggerPipeline(ctx context.Context, project_slug string, branch string, tag string, parameters interface{}) (*TriggerPipelineResponse, error) {
	path := fmt.Sprintf("/project/%s/pipeline", project_slug)
	request := TriggerPipelineRequest{
		Branch:     branch,
		Tag:        tag,
		Parameters: parameters,
	}
	result := &TriggerPipelineResponse{}
	err := s.client.http.Post(ctx, path, request, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) GetYourPipelines(ctx context.Context, project_slug string) (*CircleciPipelineResponse, error) {
	path := fmt.Sprintf("/project/%s/pipeline/mine", project_slug)
	result := &CircleciPipelineResponse{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PipelineService) GetPipelineByNumber(ctx context.Context, project_slug string, number string) (*CircleciPipeline, error) {
	path := fmt.Sprintf("/project/%s/pipeline/%s", project_slug, number)
	result := &CircleciPipeline{}
	err := s.client.http.Get(ctx, path, result)

	if err != nil {
		return nil, err
	}
	return result, nil
}

type TriggerPipelineRequest struct {
	Branch     string      `json:"branch"`
	Tag        string      `json:"tag"`
	Parameters interface{} `json:"parameters"`
}

type TriggerPipelineResponse struct {
	ID        string  `json:"id"`
	State     string  `json:"state"`
	Number    string  `json:"number"`
	CreatedAt string  `json:"created_at"`
	Message   *string `json:"message"`
}

type PipelineConfig struct {
	Source              string  `json:"source"`
	Compiled            string  `json:"compiled"`
	SetupConfig         string  `json:"setup-config"`
	CompiledSetupConfig string  `json:"compiled-setup-config"`
	Message             *string `json:"message"`
}

type ContinuePipelineRequest struct {
	ContinuationKey string      `json:"continuation-key"`
	Configuration   string      `json:"configuration"`
	Parameters      interface{} `json:"parameters"`
}

type CircleciPipelineResponse struct {
	Items         []CircleciPipeline `json:"items"`
	NextPageToken string             `json:"next_page_token"`
	Message       *string            `json:"message"`
}

type CircleciPipeline struct {
	ID                string            `json:"id"`
	Errors            []CircleciError   `json:"errors"`
	ProjectSlug       string            `json:"project_slug"`
	UpdatedAt         string            `json:"updated_at"`
	Number            int               `json:"number"`
	TriggerParameters map[string]string `json:"trigger_parameters"`
	State             string            `json:"state"`
	CreatedAt         string            `json:"created_at"`
	Trigger           Trigger           `json:"trigger"`
	Vcs               Vcs               `json:"vcs"`
	Message           *string           `json:"message"`
}

func (c *CircleciPipeline) GetMessage() string {
	if c == nil || c.Message == nil {
		return ""
	}

	return *c.Message
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
	ID        string `json:"id"`
	Name      string `json:"name"`
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
	Message       *string                `json:"message"`
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
