package services

type Platform struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Activate     bool              `json:"activate"`
	Url          string            `json:"url"`
	RestEndpoint string            `json:"rest_endpoint"`
	Property     map[string]string `json:"property"`
	Webhooks     []Webhook         `json:"webhooks"`
}

func NewPlatform(name string, url string, rest string, property map[string]string) *Platform {
	return &Platform{
		//may be need an uuid
		Id:           "",
		Name:         name,
		Activate:     true,
		Url:          url,
		RestEndpoint: rest,
		Property:     property,
		Webhooks:     []Webhook{},
	}
}

func (w *Platform) Enable() *Platform {
	w.Activate = true
	return w
}

func (w *Platform) Disable() *Platform {
	w.Activate = false
	return w
}

func (w *Platform) ChangeName(name string) *Platform {
	w.Name = name
	return w
}

func (w *Platform) ChangeUrl(url string) *Platform {
	w.Url = url
	return w
}

func (w *Platform) UpdateProperty(property map[string]string) *Platform {
	w.Property = property
	return w
}

func (w *Platform) ChangeWebhook(hook Webhook) *Platform {
	f := false
	for i := 0; i < len(w.Webhooks); i++ {
		if w.Webhooks[i].Name == hook.Name {
			w.Webhooks[i] = hook
			f = true
			break
		}
	}

	if !f {
		w.Webhooks = append(w.Webhooks, hook)
	}
	return w
}

type WebhookState interface {
	privateWebhookState()
	String() string
}

type webhookState string

func (c webhookState) privateWebhookState() {}
func (c webhookState) String() string {
	return string(c)
}

const Init webhookState = "Init"
const Creating webhookState = "Creating"
const Ready webhookState = "Ready"

type Webhook struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Activate bool              `json:"activate"`
	State    WebhookState      `json:"state"`
	Property map[string]string `json:"property"`
}

func NewWebhook(name string, url string, activate bool, state WebhookState, property map[string]string) *Webhook {
	return &Webhook{
		Name:     name,
		Url:      url,
		Activate: activate,
		State:    state,
		Property: property,
	}
}

type PlatformService struct {
	// repo
}

func NewPlatformService() *PlatformService {
	return &PlatformService{}
}

func (s *PlatformService) CreateOrUpdate(platform Platform) {
	//TODO crypt token
}

func (s *PlatformService) Get() []Platform {
	result := make([]Platform, 0)
	//TODO decrypt token
	return result
}

func (s *PlatformService) GetyIds(ids []string) []Platform {
	result := make([]Platform, 0)
	//TODO decrypt token
	return result
}
