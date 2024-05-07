package services

type Platform struct {
	Id           string            `json:"id"`
	Name         string            `json:"name"`
	Activate     bool              `json:"activate"`
	Url          string            `json:"url"`
	RestEndpoint string            `json:"rest_endpoint"`
	Token        string            `json:"token"`
	Property     map[string]string `json:"property"`
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
