package qstash

import (
	"context"
	"fmt"
	"strings"
)

var validProtocol = []string{"http://", "https://"}

// Mixing url groups with single urls is not supported.
type MessageService service

func (s *MessageService) Publish(ctx context.Context, request PublishRequest) (*PublishResponse, error) {
	valid := false
	for _, v := range validProtocol {
		if strings.HasSuffix(request.Destination, v) {
			valid = true
			break
		}
	}

	if !valid {
		return nil, fmt.Errorf("destination MUST start with 'http://' or 'https://'")
	}

	path := fmt.Sprintf("/publish/%s", request.Destination)
	result := &PublishResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) Enqueue(ctx context.Context, request EnqueueRequest) (*EnqueueResponse, error) {
	valid := false
	for _, v := range validProtocol {
		if strings.HasSuffix(request.Destination, v) {
			valid = true
			break
		}
	}

	if !valid {
		return nil, fmt.Errorf("destination MUST start with 'http://' or 'https://'")
	}

	path := fmt.Sprintf("/enqueue/%s/%s", request.QueueName, request.Destination)
	result := &EnqueueResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) Batch(ctx context.Context, request BatchRequest) (*BatchResponse, error) {
	valid := false
	for _, vv := range request {
		for _, v := range validProtocol {
			if !strings.HasSuffix(vv.Destination, v) {
				valid = true
				break
			}
		}
		if !valid {
			return nil, fmt.Errorf("destination MUST start with 'http://' or 'https://'")
		}
	}

	path := "/batch"
	result := &BatchResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) GetMessage(ctx context.Context, messageId string) (*MessageResponse, error) {
	path := fmt.Sprintf("/messages/%s", messageId)
	result := &MessageResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *MessageService) CancelMessage(ctx context.Context, messageId string) error {
	path := fmt.Sprintf("/messages/%s", messageId)
	result := ""
	return s.client.http.Delete(ctx, path, nil, &result)
}

func (s *MessageService) CancelBatchMessage(ctx context.Context, request CancelBatchRequest) (*CancelBatchResponse, error) {
	path := "/messages"
	result := &CancelBatchResponse{}
	if err := s.client.http.Delete(ctx, path, request, result); err != nil {
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

type MessageResponse struct {
	MessageId       string              `json:"messageId"`
	TopicName       string              `json:"topicName"`
	EndpointName    string              `json:"endpointName"`
	Url             string              `json:"url"`
	Method          string              `json:"method"`
	Header          map[string][]string `json:"header"`
	Body            string              `json:"body"`
	BodyBase64      string              `json:"bodyBase64"`
	MaxRetries      int                 `json:"maxRetries"`
	NotBefore       int                 `json:"notBefore"`
	CreatedAt       int                 `json:"createdAt"`
	Callback        string              `json:"callback"`
	FailureCallback string              `json:"failureCallback"`
	ScheduleId      string              `json:"scheduleId"`
	CallerIP        string              `json:"callerIP"`
}

type CancelBatchRequest struct {
	MessageIds []string `json:"messageIds"`
	QueueName  string   `json:"queueName"`
	TopicName  string   `json:"topicName"`
	Url        string   `json:"url"`
	FromDate   int      `json:"fromDate"`
	ToDate     int      `json:"toDate"`
	ScheduleId string   `json:"scheduleId"`
	CallerIP   string   `json:"callerIP"`
}

type CancelBatchResponse struct {
	Cancelled int `json:"cancelled"`
}
