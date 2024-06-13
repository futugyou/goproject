package project

import (
	"time"

	"github.com/google/uuid"

	"github.com/futugyou/infr-project/core"
	"github.com/futugyou/infr-project/domain"
)

type Project struct {
	domain.Aggregate `json:"-"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	State            ProjectState      `json:"state"`
	StartDate        *time.Time        `json:"start_date"`
	EndDate          *time.Time        `json:"end_date"`
	Platforms        []ProjectPlatform `json:"platforms"`
	Designs          []ProjectDesign   `json:"designs"`
}

func (r *Project) AggregateName() string {
	return "projects"
}

// ProjectState is the interface for webhook states.
type ProjectState interface {
	privateProjectState() // Prevents external implementation
	String() string
}

// projectState is the underlying implementation for ProjectState.
type projectState string

// privateWebhookState makes webhookState implement ProjectState.
func (c projectState) privateProjectState() {}

// String makes webhookState implement WebhookState.
func (c projectState) String() string {
	return string(c)
}

// Constants for the different webhook states.
const (
	ProjectPreparing  projectState = "preparing"
	ProjectProcessing projectState = "processing"
	ProjectFinished   projectState = "finished"
)

type ProjectPlatform struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// ref platform.Platform
	PlatformId string `json:"platform_id"`
}

type ProjectDesign struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	// ref resource.Resource
	Resources []string `json:"resources"`
}

func NewProject(name string, description string, state ProjectState, start *time.Time, end *time.Time) *Project {
	return &Project{
		Aggregate: domain.Aggregate{
			Id: uuid.New().String(),
		},
		Name:        name,
		Description: description,
		State:       state,
		StartDate:   start,
		EndDate:     end,
		Platforms:   []ProjectPlatform{},
		Designs:     []ProjectDesign{},
	}
}

func GetProjectState(rType string) ProjectState {
	switch rType {
	case "preparing":
		return ProjectPreparing
	case "processing":
		return ProjectProcessing
	case "finished":
		return ProjectFinished
	default:
		return ProjectPreparing
	}
}

func (s *Project) ChangeName(name string) *Project {
	s.Name = name
	return s
}

func (s *Project) ChangeDescription(description string) *Project {
	s.Description = description
	return s
}

func (s *Project) ChangeProjectState(state ProjectState) *Project {
	s.State = state
	return s
}

func (s *Project) ChangeStartDate(start time.Time) *Project {
	s.StartDate = &start
	return s
}

func (s *Project) ChangeEndDate(end *time.Time) *Project {
	s.EndDate = end
	return s
}

func (w *Project) UpdatePlatform(platform ProjectPlatform) *Project {
	f := false
	for i := 0; i < len(w.Platforms); i++ {
		if w.Platforms[i].Name == platform.Name {
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

func (w *Project) RemovePlatform(name string) *Project {
	for i := len(w.Platforms) - 1; i >= 0; i-- {
		if w.Platforms[i].Name == name {
			w.Platforms = append(w.Platforms[:i], w.Platforms[i+1:]...)
		}
	}
	return w
}

func (w *Project) UpdateDesign(design ProjectDesign) *Project {
	f := false
	for i := 0; i < len(w.Platforms); i++ {
		if w.Designs[i].Name == design.Name {
			w.Designs[i] = design
			f = true
			break
		}
	}

	if !f {
		w.Designs = append(w.Designs, design)
	}
	return w
}

func (w *Project) RemoveDesign(name string) *Project {
	for i := len(w.Designs) - 1; i >= 0; i-- {
		if w.Designs[i].Name == name {
			w.Designs = append(w.Designs[:i], w.Designs[i+1:]...)
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
