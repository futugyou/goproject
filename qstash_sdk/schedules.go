package qstash

import (
	"context"
	"fmt"
)

type SchedulesService service

func (s *SchedulesService) CreateSchedule(ctx context.Context, request CreateScheduleRequest) (*CreateScheduleResponse, error) {
	path := fmt.Sprintf("/v2/schedules/%s", request.Destination)
	result := &CreateScheduleResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SchedulesService) GetSchedule(ctx context.Context, scheduleId string) (*Schedule, error) {
	path := fmt.Sprintf("/v2/schedules/%s", scheduleId)
	result := &Schedule{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SchedulesService) ListSchedule(ctx context.Context) (*ScheduleList, error) {
	path := "/v2/schedules"
	result := &ScheduleList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SchedulesService) DeleteSchedule(ctx context.Context, scheduleId string) error {
	path := fmt.Sprintf("/v2/schedules/%s", scheduleId)
	result := ""
	return s.client.http.Delete(ctx, path, nil, &result)
}

func (s *SchedulesService) PauseQueue(ctx context.Context, scheduleId string) error {
	path := fmt.Sprintf("/v2/schedules/%s/pause", scheduleId)
	result := ""
	return s.client.http.Post(ctx, path, nil, &result)
}

func (s *SchedulesService) ResumeQueue(ctx context.Context, scheduleId string) error {
	path := fmt.Sprintf("/v2/schedules/%s/resume", scheduleId)
	result := ""
	return s.client.http.Post(ctx, path, nil, &result)
}

type CreateScheduleRequest struct {
	*QstashHeader `json:"-"`
	// Destination can either be a topic name or id that you configured in the Upstash console,
	// a valid url where the message gets sent to, or a valid QStash API name like api/llm.
	// If the destination is a URL, make sure the URL is prefixed with a valid protoco
	Destination string `json:"-"`
	// The raw request message passed to the endpoints as is
	Body string
}

func (r CreateScheduleRequest) GetPayload() string {
	return r.Body
}

type CreateScheduleResponse struct {
	ScheduleId string `json:"scheduleId"`
}

type Schedule struct {
	ScheduleId  string              `json:"scheduleId"`
	CreatedAt   int                 `json:"createdAt"`
	Cron        string              `json:"cron"`
	Destination string              `json:"destination"`
	Method      string              `json:"method"`
	Header      map[string][]string `json:"header"`
	Body        string              `json:"body"`
	Retries     int                 `json:"retries"`
}

type ScheduleList []Schedule
