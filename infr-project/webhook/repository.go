package webhook

import (
	"context"

	"github.com/futugyou/infr-project/domain"
)

type WebhookLogSearch struct {
	Source             string
	EventType          string
	ProviderPlatformId string
	ProviderProjectId  string
	ProviderWebhookId  string
}

type IWebhookLogRepository interface {
	domain.IRepository[WebhookLogs]
	SearchWebhookLogs(ctx context.Context, filter WebhookLogSearch) ([]WebhookLogs, error)
}
