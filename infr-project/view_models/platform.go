package viewmodels

type Property struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreatePlatformRequest struct {
	Name       string     `json:"name" validate:"required,min=3,max=50"`
	Url        string     `json:"url" validate:"required,min=3,max=150"`
	Tags       []string   `json:"tags"`
	Properties []Property `json:"properties"`
	Secrets    []Secret   `json:"secrets"`
	Provider   string     `json:"provider" validate:"oneof=vercel github circleci other"`
}

type UpdatePlatformRequest struct {
	Name       string     `json:"name" validate:"required,min=3,max=50"`
	Url        string     `json:"url" validate:"required,min=3,max=150"`
	Properties []Property `json:"properties,omitempty"`
	Secrets    []Secret   `json:"secrets"`
	Tags       []string   `json:"tags"`
	Activate   *bool      `json:"activate,omitempty"`
	Provider   string     `json:"provider" validate:"oneof=vercel github circleci other"`
}

type UpdatePlatformProjectRequest struct {
	Name       string            `json:"name" validate:"required,min=3,max=50"`
	Url        string            `json:"url" validate:"required,min=3,max=150"`
	Properties map[string]string `json:"properties,omitempty"`
}

type UpdatePlatformWebhookRequest struct {
	Name       string            `json:"name" validate:"required,min=3,max=50"`
	Url        string            `json:"url" validate:"required,min=3,max=150"`
	Activate   bool              `json:"activate" validate:"required"`
	State      string            `json:"state" validate:"oneof=Init Creating Ready"`
	Properties map[string]string `json:"properties"`
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
	Id         string     `json:"id"`
	Name       string     `json:"name"`
	Url        string     `json:"url"`
	Properties []Property `json:"properties"`
	Secrets    []Secret   `json:"secrets"`
	Webhooks   []Webhook  `json:"webhooks"`
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
