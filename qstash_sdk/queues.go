package qstash

import (
	"context"
	"fmt"
)

type QueuesService service

func (s *QueuesService) UpsertQueue(ctx context.Context, request UpsertQueueRequest) error {
	path := "/v2/queues"
	result := ""
	return s.client.http.Post(ctx, path, request, &result)
}

func (s *QueuesService) RemoveQueue(ctx context.Context, queueName string) error {
	path := fmt.Sprintf("/v2/queues/%s", queueName)
	result := ""
	return s.client.http.Delete(ctx, path, nil, &result)
}

func (s *QueuesService) ListQueue(ctx context.Context) (*QueueList, error) {
	path := "/v2/queues"
	result := &QueueList{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *QueuesService) GetQueue(ctx context.Context, queueName string) (*Queue, error) {
	path := fmt.Sprintf("/v2/queues/%s", queueName)
	result := &Queue{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *QueuesService) PauseQueue(ctx context.Context, queueName string) error {
	path := fmt.Sprintf("/v2/queues/%s/pause", queueName)
	result := ""
	return s.client.http.Post(ctx, path, nil, &result)
}

func (s *QueuesService) ResumeQueue(ctx context.Context, queueName string) error {
	path := fmt.Sprintf("/v2/queues/%s/resume", queueName)
	result := ""
	return s.client.http.Post(ctx, path, nil, &result)
}

type UpsertQueueRequest struct {
	QueueName   string `json:"queueName"`
	Parallelism int    `json:"parallelism"`
}

type QueueList []Queue

type Queue struct {
	Name        string `json:"name"`
	Parallelism int    `json:"parallelism"`
	Lag         int    `json:"lag"`
	CreatedAt   int    `json:"createdAt"`
	UpdatedAt   int    `json:"updatedAt"`
}
