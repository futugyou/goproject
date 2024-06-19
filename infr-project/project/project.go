package project

import (
	"time"

	"github.com/google/uuid"

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
	Tags             []string          `json:"tags"`
}

func (r *Project) AggregateName() string {
	return "projects"
}

func NewProject(name string, description string, state ProjectState, start *time.Time, end *time.Time, tags []string) *Project {
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
		Tags:        tags,
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

func (s *Project) ChangeTags(tags []string) *Project {
	s.Tags = tags
	return s
}

func (w *Project) UpdatePlatform(platforms []ProjectPlatform) *Project {
	w.Platforms = platforms
	return w
}

func (w *Project) UpdateDesign(designs []ProjectDesign) *Project {
	w.Designs = designs
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
