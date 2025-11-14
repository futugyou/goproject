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
	return q.dispatch(ctx, buildQStashBatch(q.destination, events))
}

func (q *QStashEventDispatcher) DispatchIntegrationEvents(ctx context.Context, events []infrastructure.Event) error {
	return q.dispatch(ctx, buildQStashBatch(q.destination, events))
}

// QStash natively supports bulk sending.
func (q *QStashEventDispatcher) dispatch(ctx context.Context, req qstash.BatchRequest) error {
	if len(req) == 0 {
		return nil
	}
	_, err := q.client.Message.Batch(ctx, req)
	return err
}

func buildQStashBatch[T interface {
	EventType() string
}](destinationFormat string, events []T) qstash.BatchRequest {
	req := qstash.BatchRequest{}
	for _, event := range events {
		bodyBytes, err := json.Marshal(event)
		if err != nil {
			continue
		}
		req = append(req, qstash.BatchRequestItem{
			Destination: fmt.Sprintf(destinationFormat, event.EventType()),
			Body:        string(bodyBytes),
		})
	}
	return req
}
