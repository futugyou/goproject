package qstash

import (
	"context"
	"fmt"
)

type MessageService service

func (s *MessageService) PublishMessage(ctx context.Context, request PublishMessageRequest) (*PublishMessageResponse, error) {
	path := fmt.Sprintf("/publish/%s", request.Destination)
	result := &PublishMessageResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

type PublishMessageRequest struct {
	*QstashHeader `json:"-"`
	// Destination can either be a topic name or id that you configured in the Upstash console,
	// a valid url where the message gets sent to, or a valid QStash API name like api/llm.
	// If the destination is a URL, make sure the URL is prefixed with a valid protoco
	Destination string `json:"-"`
	// The raw request message passed to the endpoints as is
	Body string
}

func (r PublishMessageRequest) GetPayload() string {
	return r.Body
}

type PublishMessageResponse struct {
	MessageId    string `json:"messageId"`
	Url          string `json:"url"`
	Deduplicated bool   `json:"deduplicated"`
}
