package infrastructure_qstash

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/futugyou/qstash"

	"github.com/futugyou/infr-project/domain"
)

type QStashEventPulisher struct {
	client *qstash.QstashClient
}

func NewQStashEventPulisher(client *qstash.QstashClient) *QStashEventPulisher {
	return &QStashEventPulisher{
		client: client,
	}
}

func (q *QStashEventPulisher) Publish(ctx context.Context, events []domain.IDomainEvent) error {
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

func (q *QStashEventPulisher) PublishCommon(ctx context.Context, event any, event_type string) error {
	var bodyBytes []byte
	var err error
	if bodyBytes, err = json.Marshal(event); err != nil {
		return err
	}

	qstashRequest := qstash.PublishRequest{
		Destination: fmt.Sprintf(os.Getenv("QSTASH_DESTINATION"), event_type),
		Body:        string(bodyBytes),
	}

	_, err = q.client.Message.Publish(ctx, qstashRequest)
	return err
}
