package qstash

import (
	"context"
	"fmt"
)

type WorkflowService service

func (s *WorkflowService) NotifyWorkflow(ctx context.Context, request NotifyRequest) (*NotifyResponse, error) {
	path := fmt.Sprintf("/v2/notify/%s", request.EventId)
	result := &NotifyResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

type NotifyRequest struct {
	*QstashHeader `json:"-"`
	EventId       string `json:"-"`
	Body          string
}

func (r NotifyRequest) GetPayload() string {
	return r.Body
}

type NotifyResponse []Notify

type Notify struct {
	Waiter    NotifyWaiter `json:"waiter"`
	MessageId string       `json:"messageId"`
}

type NotifyWaiter struct {
	Url           string              `json:"url"`
	Deadline      int                 `json:"deadline"`
	Headers       map[string][]string `json:"headers"`
	TimeoutUrl    string              `json:"timeoutUrl"`
	TimeoutBody   interface{}         `json:"timeoutBody"`
	TimeoutHeader map[string][]string `json:"timeoutHeaders"`
}

func (s *WorkflowService) CancelWorkflow(ctx context.Context, workflowRunId string) error {
	path := fmt.Sprintf("/v2/workflows/runs/%s", workflowRunId)
	result := ""
	return s.client.http.Post(ctx, path, nil, &result)
}
func (s *WorkflowService) BulkCancelWorkflow(ctx context.Context, request BulkCancelWorkflowRequest) (*BulkCancelWorkflowResponse, error) {
	path := "/v2/workflows/runs"
	result := &BulkCancelWorkflowResponse{}
	if err := s.client.http.Post(ctx, path, nil, result); err != nil {
		return nil, err
	}
	return result, nil
}

type BulkCancelWorkflowRequest struct {
	WorkflowRunIds []string `json:"workflowRunIds"`
	WorkflowUrl    string   `json:"workflowUrl"`
}

type BulkCancelWorkflowResponse struct {
	Cancelled int `json:"cancelled"`
}
