package webhook

import "github.com/futugyou/infr-project/domain"

type WebhookLogs struct {
	domain.Aggregate   `bson:"-"`
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
