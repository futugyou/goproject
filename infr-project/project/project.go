package project

import (
	"time"

	"github.com/google/uuid"

	"github.com/futugyou/infr-project/domain"
)

type Project struct {
	domain.Aggregate
	Name            string
	Description     string
	State           ProjectState
	StartDate       *time.Time
	EndDate         *time.Time
	Platforms       []ProjectPlatform
	Designs         []ProjectDesign
	Tags            []string
	TechStack       []string       `json:"tech_stack,omitempty"`
	CoverImageURL   string         `json:"cover_image_url,omitempty"`
	LogoURL         string         `json:"logo_url,omitempty"`
	TopologyDiagram *Diagram       `json:"topology_diagram,omitempty"`
	Milestones      []Milestone    `json:"milestones,omitempty"`
	ActivityStats   *ActivityStats `json:"activity_stats,omitempty"`
	DefaultVersion  string         `json:"default_version"`
}

func (r Project) AggregateName() string {
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
