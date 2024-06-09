package application

import (
	"time"
)

type ResourceQueryService struct {
}

func NewResourceQueryService() *ResourceQueryService {
	return &ResourceQueryService{}
}

// TODO: This is temporary because CQRS is not ready yet
func (s *ResourceQueryService) GetAllResourceSnapshots() ([]ResourceDetail, error) {
	return []ResourceDetail{}, nil
}

type ResourceDetail struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Data      string    `json:"data"`
	IsDelete  bool      `json:"is_deleted"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
