package qstashdispatcherimpl

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/futugyou/qstash"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyou/domaincore/infrastructure"
)

type QStashEventDispatcher struct {
	client      *qstash.QstashClient
	destination string
}

func NewQStashEventDispatcher(token string, destination string) *QStashEventDispatcher {
	qstashClient := qstash.NewClient(token)
	return &QStashEventDispatcher{
		client:      qstashClient,
		destination: destination,
	}
}

func (q *QStashEventDispatcher) DispatchDomainEvents(ctx context.Context, events []domain.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	qstashRequest := qstash.BatchRequest{}
	for _, event := range events {
		if bodyBytes, err := json.Marshal(event); err == nil {
			qstashRequest = append(qstashRequest, qstash.BatchRequestItem{
				Destination: fmt.Sprintf(q.destination, event.EventType()),
				Body:        string(bodyBytes),
			})
		}
	}

	if len(qstashRequest) == 0 {
		return nil
	}

	_, err := q.client.Message.Batch(ctx, qstashRequest)
	return err
}

func (q *QStashEventDispatcher) DispatchIntegrationEvent(ctx context.Context, event infrastructure.Event) error {
	var bodyBytes []byte
	var err error
	if bodyBytes, err = json.Marshal(event); err != nil {
		return err
	}

	qstashRequest := qstash.PublishRequest{
		Destination: fmt.Sprintf(q.destination, event.EventType()),
		Body:        string(bodyBytes),
	}

	_, err = q.client.Message.Publish(ctx, qstashRequest)
	return err
}
