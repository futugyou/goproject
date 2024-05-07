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

type State interface {
	privateState()
	String() string
}

type state string

func (c state) privateState() {}
func (c state) String() string {
	return string(c)
}

const Init state = "Init"
const Creating state = "Creating"
const Ready state = "Ready"

type Webhook struct {
	Name     string            `json:"name"`
	Url      string            `json:"url"`
	Activate bool              `json:"activate"`
	State    State             `json:"state"`
	Property map[string]string `json:"property"`
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
