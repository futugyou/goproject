package services

import (
	"time"

	"github.com/futugyou/infr-project/core"
)

type Project struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	ProjectDate time.Time  `json:"project_date"`
	Platforms   []Platform `json:"platforms"`
}

type ProjectService struct {
}

func NewProjectService() *ProjectService {
	return &ProjectService{}
}

func (s *ProjectService) Get(filter core.Search) []Project {
	result := make([]Project, 0)
	return result
}
