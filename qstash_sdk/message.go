package qstash

import (
	"context"
	"fmt"
)

// Only single URL is supported, not URL groups.
type MessageService service

func (s *MessageService) Publish(ctx context.Context, request PublishRequest) (*PublishResponse, error) {
	path := fmt.Sprintf("/publish/%s", request.Destination)
	result := &PublishResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) Enqueue(ctx context.Context, request EnqueueRequest) (*EnqueueResponse, error) {
	path := fmt.Sprintf("/enqueue/%s/%s", request.QueueName, request.Destination)
	result := &EnqueueResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) Batch(ctx context.Context, request []BatchRequest) (*BatchResponse, error) {
	path := "/batch"
	result := &BatchResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

type EnqueueRequest struct {
	*QstashHeader `json:"-"`
	// Destination can either be a topic name or id that you configured in the Upstash console,
	// a valid url where the message gets sent to, or a valid QStash API name like api/llm.
	// If the destination is a URL, make sure the URL is prefixed with a valid protoco
	Destination string `json:"-"`
	QueueName   string `json:"-"`
	// The raw request message passed to the endpoints as is
	Body string
}

func (r EnqueueRequest) GetPayload() string {
	return r.Body
}

type PublishRequest struct {
	*QstashHeader `json:"-"`
	// Destination can either be a topic name or id that you configured in the Upstash console,
	// a valid url where the message gets sent to, or a valid QStash API name like api/llm.
	// If the destination is a URL, make sure the URL is prefixed with a valid protoco
	Destination string `json:"-"`
	// The raw request message passed to the endpoints as is
	Body string
}

func (r PublishRequest) GetPayload() string {
	return r.Body
}

type PublishResponse struct {
	MessageId    string `json:"messageId"`
	Url          string `json:"url"`
	Deduplicated bool   `json:"deduplicated"`
}

type EnqueueResponse struct {
	MessageId    string `json:"messageId"`
	Url          string `json:"url"`
	Deduplicated bool   `json:"deduplicated"`
}

type BatchRequest []BatchRequestItem

type BatchRequestItem struct {
	Destination string            `json:"destination"`
	Headers     map[string]string `json:"headers"`
	Queue       string            `json:"queue"`
}

type BatchResponse []BatchResponseItem

type BatchResponseItem struct {
	MessageId    string `json:"messageId"`
	Url          string `json:"url"`
	Deduplicated bool   `json:"deduplicated"`
}
