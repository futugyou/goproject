package viewmodel

type UpdatePlatformProjectRequest struct {
	Name              string     `json:"name" validate:"required,min=3,max=50"`
	Url               string     `json:"url" validate:"required,min=3,max=150"`
	Secrets           []Secret   `json:"secrets" validate:"required"` // only Key and VaultId in request
	Description       string     `json:"description" validate:"required,min=3,max=250"`
	Properties        []Property `json:"properties" validate:"required"`
	Tags              []string   `json:"tags"`
	ProviderProjectID string     `json:"provider_project_id"`
	Operate           string     `json:"operate" validate:"oneof=upsert sync"`
	ImportWebhooks    bool       `json:"import_webhooks"`
}

type PlatformProject struct {
	ID                string                   `json:"id"`
	Name              string                   `json:"name"`
	Url               string                   `json:"url"`
	ImageUrl          string                   `json:"image_url"`
	Description       string                   `json:"description"`
	Properties        []Property               `json:"properties"`
	Secrets           []Secret                 `json:"secrets"`
	Webhook           *Webhook                 `json:"webhook"`
	ProviderProjectID string                   `json:"provider_project_id"`
	Followed          bool                     `json:"followed"`
	ProviderProject   *PlatformProviderProject `json:"provider_project"`
}

type PlatformProviderProject struct {
	ID                   string                `json:"id"`
	Name                 string                `json:"name"`
	Url                  string                `json:"url"`
	Description          string                `json:"description"`
	WebHook              *Webhook              `json:"webhook"`
	Properties           []Property            `json:"properties"`
	EnvironmentVariables []EnvironmentVariable `json:"environment_variables"`
	Environments         []string              `json:"environments"`
	Workflows            []Workflow            `json:"workflows"`
	WorkflowRuns         []WorkflowRun         `json:"workflow_runs"`
	Deployments          []Deployment          `json:"deployments"`
	BadgeURL             string                `json:"badge_url"`
	BadgeMarkdown        string                `json:"badge_markdown"`
	Tags                 []string              `json:"tags"`
	Readme               string                `json:"readme"`
}

type EnvironmentVariable struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type Workflow struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	BadgeURL      string `json:"badge_url"`
	BadgeMarkdown string `json:"badge_markdown"`
}

type WorkflowRun struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	CreatedAt     string `json:"createdAt"`
	BadgeURL      string `json:"badge_url"`
	BadgeMarkdown string `json:"badge_markdown"`
}

type Deployment struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Environment   string `json:"environment"`
	ReadyState    string `json:"readyState"`
	ReadySubstate string `json:"readySubstate"`
	CreatedAt     string `json:"createdAt"`
	BadgeURL      string `json:"badge_url"`
	BadgeMarkdown string `json:"badge_markdown"`
	Description   string `json:"description"`
}

type Webhook struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Url        string     `json:"url"`
	Events     []string   `json:"events"`
	Activate   bool       `json:"activate"`
	State      string     `json:"state"`
	Properties []Property `json:"properties"`
	Secrets    []Secret   `json:"secrets"`
	Followed   bool       `json:"followed"`
}
