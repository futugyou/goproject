package viewmodels

import "time"

type CreateProjectRequest struct {
	Name         string     `json:"name" validate:"required,min=3,max=50"`
	Description  string     `json:"description" validate:"required,min=3,max=500"`
	ProjectState *string    `json:"state" validate:"oneof=preparing processing finished"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
}

type UpdateProjectRequest struct {
	Name         *string    `json:"name" validate:"min=3,max=50"`
	Description  *string    `json:"description" validate:"min=3,max=500"`
	ProjectState *string    `json:"state" validate:"oneof=preparing processing finished"`
	StartTime    *time.Time `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
}

type UpdateProjectPlatformRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50"`
	Description string `json:"description" validate:"min=3,max=500"`
	PlatformId  string `json:"platform_id" validate:"required,min=3,max=50"`
}
