package webhook

import (
	"time"

	tool "github.com/futugyou/extensions"
	"github.com/google/uuid"

	"github.com/futugyou/infr-project/domain"
)

type WebhookLogs struct {
	domain.Aggregate   `bson:",inline"`
	Source             string `bson:"source"` // github/vercel/circleci
	EventType          string `bson:"event_type"`
	ProviderPlatformId string `bson:"provider_platform_id"`
	ProviderProjectId  string `bson:"provider_project_id"`
	ProviderWebhookId  string `bson:"provider_webhook_id"`
	Data               string `bson:"data"`
	HappenedAt         string `bson:"happened_at"`
}

func (r WebhookLogs) AggregateName() string {
	return "platform_webhook_logs"
}

func (r WebhookLogs) Verify(secret string, signature string) (bool, error) {
	return tool.VerifySignatureHMAC(secret, signature, r.Data)
}

func NewWebhookLogs(source, eventType, providerPlatformId, providerProjectId, providerWebhookId, data string) *WebhookLogs {
	return &WebhookLogs{
		Aggregate: domain.Aggregate{
			Id: uuid.NewString(),
		},
		Source:             source,
		EventType:          eventType,
		ProviderPlatformId: providerPlatformId,
		ProviderProjectId:  providerProjectId,
		ProviderWebhookId:  providerWebhookId,
		Data:               data,
		HappenedAt:         time.Now().UTC().Format(time.RFC3339Nano),
	}
}
