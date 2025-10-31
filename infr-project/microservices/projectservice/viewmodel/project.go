package viewmodel

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
	ProjectID   string `json:"project_id" validate:"required,min=3,max=50"`
	PlatformID  string `json:"platform_id" validate:"required,min=3,max=50"`
}

type UpdateProjectDesignRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ResourceID      string `json:"resource_id"`
	ResourceVersion int    `json:"resource_version"`
}

type ProjectView struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	State       string            `json:"state"`
	StartDate   *time.Time        `json:"start_time"`
	EndDate     *time.Time        `json:"end_time"`
	Platforms   []ProjectPlatform `json:"platforms"`
	Designs     []ProjectDesign   `json:"designs"`
	Tags        []string          `json:"tags"`
}

type ProjectPlatform struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PlatformID  string `json:"platform_id"`
	ProjectID   string `json:"project_id"`
}

type ProjectDesign struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	ResourceID      string `json:"resource_id"`
	ResourceVersion int    `json:"resource_version"`
}
