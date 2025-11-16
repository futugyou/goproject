package domain

import (
	"time"

	tool "github.com/futugyou/extensions"
	"github.com/google/uuid"

	"github.com/futugyou/domaincore/domain"
)

type WebhookLogs struct {
	domain.Aggregate   `bson:",inline"`
	Source             string `bson:"source"` // github/vercel/circleci
	EventType          string `bson:"event_type"`
	ProviderPlatformId string `bson:"provider_platform_id"`
	ProviderProjectId  string `bson:"provider_project_id"`
	ProviderWebhookId  string `bson:"provider_webhook_id"`
	Signature          string `bson:"signature"`
	Data               string `bson:"data"`
	HappenedAt         string `bson:"happened_at"`
}

func (r WebhookLogs) Verify(secret string) (bool, error) {
	return tool.VerifySignatureHMAC(secret, r.Signature, r.Data)
}

func NewWebhookLogs(source, eventType, providerPlatformId, providerProjectId, providerWebhookId, data string, signature string) *WebhookLogs {
	return &WebhookLogs{
		Aggregate: domain.Aggregate{
			ID: uuid.NewString(),
		},
		Source:             source,
		EventType:          eventType,
		ProviderPlatformId: providerPlatformId,
		ProviderProjectId:  providerProjectId,
		ProviderWebhookId:  providerWebhookId,
		Data:               data,
		HappenedAt:         time.Now().UTC().Format(time.RFC3339Nano),
		Signature:          signature,
	}
}
