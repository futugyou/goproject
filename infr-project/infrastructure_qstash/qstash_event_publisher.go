package infrastructure_qstash

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/futugyou/qstash"

	"github.com/futugyou/infr-project/domain"
)

type QStashEventPulisher[Event domain.IDomainEvent] struct {
	client *qstash.QstashClient
}

func NewQStashEventPulisher[Event domain.IDomainEvent](client *qstash.QstashClient) *QStashEventPulisher[Event] {
	return &QStashEventPulisher[Event]{
		client: client,
	}
}

func (q *QStashEventPulisher[Event]) Publish(ctx context.Context, events []Event) error {
	if len(events) == 0 {
		return nil
	}

	qstashRequest := qstash.BatchRequest{}
	for _, event := range events {
		if bodyBytes, err := json.Marshal(event); err == nil {
			qstashRequest = append(qstashRequest, qstash.BatchRequestItem{
				Destination: fmt.Sprintf(os.Getenv("QSTASH_DESTINATION"), event.EventType()),
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
