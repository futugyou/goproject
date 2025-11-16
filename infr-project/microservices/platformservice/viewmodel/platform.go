package viewmodel

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

type PlatformView struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Activate  bool     `json:"activate"`
	Url       string   `json:"url"`
	Tags      []string `json:"tags"`
	IsDeleted bool     `json:"is_deleted"`
	Provider  string   `json:"provider"`
}

type PlatformDetailView struct {
	ID         string            `json:"id"`
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
	VaultID   string `json:"vault_id"`
	VaultKey  string `json:"vault_key,omitempty"`
	MaskValue string `json:"mask_value,omitempty"`
}

type SearchPlatformsRequest struct {
	Name     string   `json:"name"`
	Activate *bool    `json:"activate"`
	Tags     []string `json:"tags"`
	Page     int      `json:"page"`
	Size     int      `json:"size"`
}
