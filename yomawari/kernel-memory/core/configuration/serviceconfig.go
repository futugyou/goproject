package configuration

import "math"

type ServiceConfig struct {
	RunWebService      bool
	RunHandlers        bool
	OpenApiEnabled     bool
	SendSSEDoneMessage bool
	Handlers           map[string]HandlerConfig
	MaxUploadSizeMb    *int64
}

func (s *ServiceConfig) GetMaxUploadSizeInBytes() *int64 {
	if s.MaxUploadSizeMb != nil {
		r := int64(math.Max(1, float64(*s.MaxUploadSizeMb))) * 1024 * 1024
		return &r
	}
	return nil
}
