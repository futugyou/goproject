package viewmodels

import "time"

type CreateProjectRequest struct {
	Name         string     `json:"name" validate:"required,min=3,max=50"`
	Description  string     `json:"description" validate:"required,min=3,max=500"`
	ProjectState string     `json:"state" validate:"oneof=preparing processing finished"`
	Tags         []string   `json:"tags"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
}

type UpdateProjectRequest struct {
	Name         string     `json:"name" validate:"min=3,max=50"`
	Description  string     `json:"description" validate:"min=3,max=500"`
	ProjectState string     `json:"state" validate:"oneof=preparing processing finished"`
	Tags         []string   `json:"tags"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
}

type UpdateProjectPlatformRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"min=3,max=500"`
	ProjectId   string `json:"project_id" validate:"required,min=3,max=50"`
}

type UpdateProjectDesignRequest struct {
	Name        string   `json:"name" validate:"required,min=3,max=50"`
	Description string   `json:"description" validate:"min=3,max=500"`
	Resources   []string `json:"resources" validate:"required"`
}

type ProjectView struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	State       string            `json:"state"`
	StartDate   *time.Time        `json:"start_date"`
	EndDate     *time.Time        `json:"end_date"`
	Platforms   []ProjectPlatform `json:"platforms"`
	Designs     []ProjectDesign   `json:"designs"`
	Tags        []string          `json:"tags"`
}

type ProjectPlatform struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ProjectId   string `json:"project_id"`
}

type ProjectDesign struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Resources   []string `json:"resources"`
}
