package abstractions

const (
	ModelIdKey    string = "ModelId"
	EndpointKey   string = "Endpoint"
	ApiVersionKey string = "ApiVersion"
)

type IAIService interface {
	GetAttributes() map[string]interface{}
	GetModelId() string
	GetEndpoint() string
	GetApiVersion() string
}

type DefaultAIService struct {
	Attributes map[string]interface{}
}

func NewDefaultAIService(attributes map[string]interface{}) *DefaultAIService {
	return &DefaultAIService{
		Attributes: attributes,
	}
}

func (service DefaultAIService) GetEndpoint() string {
	return service.getAttribute(EndpointKey)
}

func (service DefaultAIService) GetModelId() string {
	return service.getAttribute(ModelIdKey)
}

func (service DefaultAIService) GetApiVersion() string {
	return service.getAttribute(ApiVersionKey)
}

func (service DefaultAIService) GetAttributes() map[string]interface{} {
	return service.Attributes
}

func (service DefaultAIService) getAttribute(key string) string {
	if value, ok := service.Attributes[key].(string); ok {
		return value
	}
	return ""
}
