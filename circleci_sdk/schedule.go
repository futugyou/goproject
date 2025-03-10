package circleci

import (
	"context"
	"fmt"
)

type ScheduleService service

func (s *ScheduleService) GetAllSchedules(ctx context.Context, project_slug string, page_token string) (*ScheduleList, error) {
	path := fmt.Sprintf("/project/%s/schedule", project_slug)
	if len(page_token) > 0 {
		path += ("?page-token=" + page_token)
	}
	result := &ScheduleList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, project_slug string, request CreateScheduleRequest) (*ScheduleInfo, error) {
	path := fmt.Sprintf("/project/%s/schedule", project_slug)
	result := &ScheduleInfo{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, schedule_id string) (*ScheduleInfo, error) {
	path := fmt.Sprintf("/schedule/%s", schedule_id)
	result := &ScheduleInfo{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, schedule_id string, request CreateScheduleRequest) (*ScheduleInfo, error) {
	path := fmt.Sprintf("/schedule/%s", schedule_id)
	result := &ScheduleInfo{}
	if err := s.client.http.Patch(ctx, path, request, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, schedule_id string) (*BaseResponse, error) {
	path := fmt.Sprintf("/schedule/%s", schedule_id)
	result := &BaseResponse{}
	if err := s.client.http.Delete(ctx, path, result); err != nil {
		return nil, err
	}
	return result, nil
}

type CreateScheduleRequest struct {
	Name             string     `json:"name"`
	Timetable        Timetable  `json:"timetable"`
	AttributionActor string     `json:"attribution-actor"` // "current" "system"
	Parameters       Parameters `json:"parameters"`
	Description      string     `json:"description"`
}

type ScheduleList struct {
	Items         []ScheduleInfo `json:"items"`
	NextPageToken string         `json:"next_page_token"`
	Message       *string        `json:"message,omitempty"`
}

type ScheduleInfo struct {
	ID          string     `json:"id"`
	Timetable   Timetable  `json:"timetable"`
	UpdatedAt   string     `json:"updated-at"`
	Name        string     `json:"name"`
	CreatedAt   string     `json:"created-at"`
	ProjectSlug string     `json:"project-slug"`
	Parameters  Parameters `json:"parameters"`
	Actor       Actor      `json:"actor"`
	Description string     `json:"description"`
	Message     *string    `json:"message,omitempty"`
}

type Parameters struct {
	DeployProd bool   `json:"deploy_prod"`
	Branch     string `json:"branch"`
}

type Timetable struct {
	PerHour     int64    `json:"per-hour"`
	HoursOfDay  []int64  `json:"hours-of-day"`
	DaysOfWeek  []string `json:"days-of-week"`
	DaysOfMonth []int64  `json:"days-of-month"`
	Months      []string `json:"months"`
}
