package circleci

type CircleciClient struct {
	http     *httpClient
	common   service
	Context  *ContextService
	Insights *InsightsService
	Job      *JobService
	Pipeline *PipelineService
	Policy   *PolicyService
	Project  *ProjectService
	Schedule *ScheduleService
	Token    *TokenService
	Usage    *UsageService
	User     *UserService
	Webhook  *WebhookService
	Workflow *WorkflowService
}

type service struct {
	client *CircleciClient
}

const circleci_url string = "https://circleci.com/api/v2"
const circleci_url_v1 string = "https://circleci.com/api/v1.1"

func NewClientV2(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: NewHttpClientWithHeader(circleci_url, header),
	}
	c.initialize()
	return c
}

func NewClientV1(token string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: NewHttpClientWithHeader(circleci_url_v1, header),
	}
	c.initialize()
	return c
}

func NewClientWithFelfHosted(token string, hostUrl string) *CircleciClient {
	header := make(map[string]string)
	header["Circle-Token"] = token

	c := &CircleciClient{
		http: NewHttpClientWithHeader(hostUrl, header),
	}
	c.initialize()
	return c
}

func (c *CircleciClient) initialize() {
	c.common.client = c
	c.Context = (*ContextService)(&c.common)
	c.Insights = (*InsightsService)(&c.common)
	c.Job = (*JobService)(&c.common)
	c.Pipeline = (*PipelineService)(&c.common)
	c.Policy = (*PolicyService)(&c.common)
	c.Project = (*ProjectService)(&c.common)
	c.Schedule = (*ScheduleService)(&c.common)
	c.Token = (*TokenService)(&c.common)
	c.Usage = (*UsageService)(&c.common)
	c.User = (*UserService)(&c.common)
	c.Webhook = (*WebhookService)(&c.common)
	c.Workflow = (*WorkflowService)(&c.common)
}

type BaseResponse struct {
	Message *string `json:"message,omitempty"`
}
