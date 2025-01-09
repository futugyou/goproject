package viewmodels

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreatePlatformRequest struct {
	Name       string     `json:"name" validate:"required,min=3,max=50"`
	Url        string     `json:"url" validate:"required,min=3,max=150"`
	Tags       []string   `json:"tags" validate:"required"`
	Properties []Property `json:"properties" validate:"required"`
	Secrets    []Secret   `json:"secrets" validate:"required"` // only Key and VaultId in request
	Provider   string     `json:"provider" validate:"oneof=vercel github circleci other"`
}

type UpdatePlatformRequest struct {
	Name       string     `json:"name" validate:"required,min=3,max=50"`
	Url        string     `json:"url" validate:"required,min=3,max=150"`
	Properties []Property `json:"properties" validate:"required"`
	Secrets    []Secret   `json:"secrets" validate:"required"` // only Key and VaultId in request
	Tags       []string   `json:"tags" validate:"required"`
	Provider   string     `json:"provider" validate:"oneof=vercel github circleci other"`
}

type UpdatePlatformProjectRequest struct {
	Name              string     `json:"name" validate:"required,min=3,max=50"`
	Url               string     `json:"url" validate:"required,min=3,max=150"`
	Secrets           []Secret   `json:"secrets" validate:"required"` // only Key and VaultId in request
	Properties        []Property `json:"properties" validate:"required"`
	ProviderProjectId string     `json:"provider_project_id"`
	Operate           string     `json:"operate" validate:"oneof=upsert sync"`
}

type UpdatePlatformWebhookRequest struct {
	Name       string     `json:"name" validate:"required,min=3,max=50"`
	Url        string     `json:"url" validate:"required,min=3,max=150"`
	Activate   bool       `json:"activate"`
	State      string     `json:"state" validate:"oneof=Init Creating Ready"`
	Secrets    []Secret   `json:"secrets" validate:"required"` // only Key and VaultId in request
	Properties []Property `json:"properties" validate:"required"`
	Sync       bool       `json:"sync"`
}

type RemoveWebhookRequest struct {
	PlatformId string `json:"-"` // Platform Id or Name
	ProjectId  string `json:"-"`
	HookName   string `json:"-"`
	Sync       bool   `json:"sync"`
}

type PlatformView struct {
	Id        string   `json:"id"`
	Name      string   `json:"name"`
	Activate  bool     `json:"activate"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	IsDeleted bool     `json:"is_deleted"`
	Provider  string   `json:"provider"`
}

type PlatformDetailView struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Activate   bool              `json:"activate"`
	Url        string            `json:"url"`
	Properties []Property        `json:"properties"`
	Secrets    []Secret          `json:"secrets"`
	Projects   []PlatformProject `json:"projects"`
	Tags       []string          `json:"tags"`
	IsDeleted  bool              `json:"is_deleted"`
	Provider   string            `json:"provider"`
}

type Secret struct {
	Key       string `json:"key"` //vault aliases
	VaultId   string `json:"vault_id"`
	VaultKey  string `json:"vault_key,omitempty"`
	MaskValue string `json:"mask_value,omitempty"`
}

type PlatformProject struct {
	Id                string       `json:"id"`
	Name              string       `json:"name"`
	Url               string       `json:"url"`
	Description       string       `json:"description"`
	Properties        []Property   `json:"properties"`
	Secrets           []Secret     `json:"secrets"`
	Webhooks          []Webhook    `json:"webhooks"`
	Followed          bool         `json:"followed"`
	ProviderProjectId string       `json:"provider_project_id"`
	Environments      []ProjectEnv `json:"environments"`
	Workflows         []Workflow   `json:"workflows"`
	Deployments       []Deployment `json:"deployments"`
	BadgeURL          string       `json:"badge_url"`
	BadgeMarkdown     string       `json:"badge_markdown"`
}

type ProjectEnv struct {
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

type Deployment struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Plan          string `json:"plan"`
	ReadyState    string `json:"readyState"`
	ReadySubstate string `json:"readySubstate"`
	CreatedAt     string `json:"createdAt"`
	BadgeURL      string `json:"badge_url"`
	BadgeMarkdown string `json:"badge_markdown"`
}

type Webhook struct {
	Name       string     `json:"name"`
	Url        string     `json:"url"`
	Activate   bool       `json:"activate"`
	State      string     `json:"state"`
	Properties []Property `json:"properties"`
	Secrets    []Secret   `json:"secrets"`
}

type SearchPlatformsRequest struct {
	Name     string   `json:"name"`
	Activate *bool    `json:"activate"`
	Tags     []string `json:"tags"`
	Page     int      `json:"page"`
	Size     int      `json:"size"`
}
