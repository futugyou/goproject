package qstash

import (
	"context"
	"fmt"
	"net/url"
)

type EventsService service

func (s *EventsService) ListEvents(ctx context.Context, query ListEventsQuery) (*QstashEventResponse, error) {
	u := &url.URL{
		Path: "/events",
	}

	params := url.Values{}
	if query.Cursor != nil {
		params.Add("cursor", *query.Cursor)
	}
	if query.MessageId != nil {
		params.Add("messageId", *query.MessageId)
	}
	if query.State != nil {
		params.Add("state", (string)(*query.State))
	}
	if query.Url != nil {
		params.Add("url", *query.Url)
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
		if *query.Count >= 1000 {
			params.Add("count", fmt.Sprintf("%d", 1000))
		} else {
			params.Add("count", fmt.Sprintf("%d", *query.Count))
		}
	}
	if query.FromDate != nil {
		params.Add("fromDate", fmt.Sprintf("%d", *query.FromDate))
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
	result := &QstashEventResponse{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

type ListEventsQuery struct {
	Cursor     *string
	MessageId  *string
	State      *EventState
	Url        *string
	TopicName  *string
	ScheduleId *string
	QueueName  *string
	FromDate   *int
	ToDate     *int
	Count      *int
	Order      *string
}

type EventState string

const ACTIVE EventState = "ACTIVE"
const DELIVERED EventState = "DELIVERED"
const RETRY EventState = "RETRY"
const CANCEL_REQUESTED EventState = "CANCEL_REQUESTED"
const CANCELLED EventState = "CANCELLED"
const FAILED EventState = "FAILED"
const CREATED EventState = "CREATED"
const ERROR EventState = "ERROR"

type QstashEventResponse struct {
	Cursor string        `json:"cursor"`
	Events []QstashEvent `json:"events"`
}

type QstashEvent struct {
	Time      string              `json:"time"`
	MessageId string              `json:"messageId"`
	State     string              `json:"state"`
	Url       string              `json:"url"`
	Header    map[string][]string `json:"header"`
	Body      string              `json:"body"`
}
