package project

import (
	"time"

	"github.com/google/uuid"

	"github.com/futugyou/infr-project/core"
	"github.com/futugyou/infr-project/domain"
	"github.com/futugyou/infr-project/platform"
	"github.com/futugyou/infr-project/resource"
)

type Project struct {
	domain.Aggregate `json:"-"`
	Name             string              `json:"name"`
	Description      string              `json:"description"`
	StartDate        time.Time           `json:"start_date"`
	EndDate          *time.Time          `json:"end_date"`
	Platforms        []platform.Platform `json:"platforms"`
	Resources        []resource.Resource `json:"resources"`
}

func (r *Project) AggregateName() string {
	return "projects"
}

func NewProject(name string, description string) *Project {
	return &Project{
		Aggregate: domain.Aggregate{
			Id: uuid.New().String(),
		},
		Name:        name,
		Description: description,
		StartDate:   time.Now().UTC(),
		Platforms:   []platform.Platform{},
	}
}

func (s *Project) ChangeName(name string) *Project {
	s.Name = name
	return s
}

func (w *Project) UpdatePlatform(platform platform.Platform) *Project {
	f := false
	for i := 0; i < len(w.Platforms); i++ {
		if w.Platforms[i].Id == platform.Id {
			w.Platforms[i] = platform
			f = true
			break
		}
	}

	if !f {
		w.Platforms = append(w.Platforms, platform)
	}
	return w
}

func (w *Project) RemovePlatform(id string) *Project {
	for i := len(w.Platforms) - 1; i >= 0; i-- {
		if w.Platforms[i].Id == id {
			w.Platforms = append(w.Platforms[:i], w.Platforms[i+1:]...)
		}
	}
	return w
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
