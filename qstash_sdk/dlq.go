package qstash

import (
	"context"
	"fmt"
	"net/url"
)

type DeadLetterQueuesService service

func (s *DeadLetterQueuesService) ListMessages(ctx context.Context, query ListMessagesQuery) (*ListDLQResponse, error) {
	u := &url.URL{
		Path: "/v2/dlq",
	}

	params := url.Values{}
	if query.Cursor != nil {
		params.Add("cursor", *query.Cursor)
	}
	if query.MessageId != nil {
		params.Add("messageId", *query.MessageId)
	}
	if query.Url != nil {
		params.Add("url", *query.Url)
	}
	if query.Api != nil {
		params.Add("api", *query.Api)
	}
	if query.CallerIp != nil {
		params.Add("callerIp", *query.CallerIp)
	}
	if query.TopicName != nil {
		params.Add("topicName", *query.TopicName)
	}
	if query.ScheduleId != nil {
		params.Add("scheduleId", *query.ScheduleId)
	}
	if query.QueueName != nil {
		params.Add("queueName", *query.QueueName)
	}
	if query.Count != nil {
		if *query.Count >= 100 {
			params.Add("count", fmt.Sprintf("%d", 100))
		} else {
			params.Add("count", fmt.Sprintf("%d", *query.Count))
		}
	}
	if query.FromDate != nil {
		params.Add("fromDate", fmt.Sprintf("%d", *query.FromDate))
	}
	if query.ResponseStatus != nil {
		params.Add("responseStatus", fmt.Sprintf("%d", *query.ResponseStatus))
	}
	if query.ToDate != nil {
		params.Add("toDate", fmt.Sprintf("%d", *query.ToDate))
	}
	if query.Order != nil {
		if *query.Order == "earliestFirst" {
			params.Add("order", "earliestFirst")
		} else {
			params.Add("order", "latestFirst")
		}
	}

	u.RawQuery = params.Encode()
	path := u.String()
	result := &ListDLQResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DeadLetterQueuesService) GetMessage(ctx context.Context, dlqId string) (*DLQMessage, error) {
	path := fmt.Sprintf("/v2/dlq/%s", dlqId)
	result := &DLQMessage{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DeadLetterQueuesService) DeleteMessage(ctx context.Context, dlqId string) error {
	path := fmt.Sprintf("/v2/dlq/%s", dlqId)
	result := ""
	return s.client.http.Delete(ctx, path, nil, &result)
}

func (s *DeadLetterQueuesService) DeleteMultipleMessage(ctx context.Context, dlqIds []string) (*MultipleDeleteResponse, error) {
	path := "/v2/dlq"
	request := struct {
		DlqIds []string `json:"dlqIds"`
	}{
		DlqIds: dlqIds,
	}

	result := MultipleDeleteResponse{}
	if err := s.client.http.Delete(ctx, path, request, result); err != nil {
		return nil, err
	}

	return &result, nil
}

type MultipleDeleteResponse struct {
	Deleted int `json:"deleted"`
}

type ListMessagesQuery struct {
	Cursor         *string
	MessageId      *string
	Url            *string
	TopicName      *string
	ScheduleId     *string
	QueueName      *string
	Api            *string
	ResponseStatus *int
	FromDate       *int
	ToDate         *int
	Count          *int
	Order          *string
	CallerIp       *string
}

type ListDLQResponse struct {
	Cursor   string       `json:"cursor"`
	Messages []DLQMessage `json:"messages"`
}

type DLQMessage struct {
	MessageId          string              `json:"messageId"`
	TopicName          string              `json:"topicName,omitempty"`
	EndpointName       string              `json:"endpointName,omitempty"`
	Url                string              `json:"url"`
	Method             string              `json:"method,omitempty"`
	Header             map[string][]string `json:"header,omitempty"`
	Body               string              `json:"body,omitempty"`
	BodyBase64         string              `json:"bodyBase64,omitempty"`
	MaxRetries         int                 `json:"maxRetries,omitempty"`
	NotBefore          int                 `json:"notBefore,omitempty"`
	CreatedAt          int                 `json:"createdAt"`
	Callback           string              `json:"callback,omitempty"`
	FailureCallback    string              `json:"failureCallback,omitempty"`
	ScheduleId         string              `json:"scheduleId,omitempty"`
	CallerIP           string              `json:"callerIP"`
	DlqId              string              `json:"dlqId"`
	QueueName          string              `json:"queueName,omitempty"`
	ResponseStatus     int                 `json:"responseStatus,omitempty"`
	ResponseHeader     map[string][]string `json:"responseHeader,omitempty"`
	ResponseBody       string              `json:"responseBody,omitempty"`
	ResponseBodyBase64 string              `json:"responseBodyBase64,omitempty"`
}
