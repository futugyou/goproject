package services

import (
	"time"

	"github.com/google/uuid"

	"github.com/futugyou/infr-project/core"
	"github.com/futugyou/infr-project/platform"
)

type Project struct {
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	ProjectDate time.Time           `json:"project_date"`
	Platforms   []platform.Platform `json:"platforms"`
}

func NewProject(name string) *Project {
	return &Project{
		Id:          uuid.New().String(),
		Name:        name,
		ProjectDate: time.Now().UTC(),
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
