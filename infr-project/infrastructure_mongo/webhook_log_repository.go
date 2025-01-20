package infrastructure_mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/futugyou/infr-project/extensions"
	"github.com/futugyou/infr-project/webhook"
)

type WebhookLogRepository struct {
	BaseRepository[webhook.WebhookLogs]
}

func NewWebhookLogRepository(client *mongo.Client, config DBConfig) *WebhookLogRepository {
	return &WebhookLogRepository{
		BaseRepository: *NewBaseRepository[webhook.WebhookLogs](client, config),
	}
}

func (r *WebhookLogRepository) SearchWebhookLogs(ctx context.Context, filter *webhook.WebhookLogSearch) ([]webhook.WebhookLogs, error) {
	f := map[string]interface{}{}
	if filter != nil {
		if filter.Source != nil {
			f["source"] = *filter.Source
		}
		if filter.EventType != nil {
			f["event_type"] = *filter.EventType
		}
		if filter.ProviderPlatformId != nil {
			f["provider_platform_id"] = *filter.ProviderPlatformId
		}
		if filter.ProviderProjectId != nil {
			f["provider_project_id"] = *filter.ProviderProjectId
		}
		if filter.ProviderWebhookId != nil {
			f["provider_webhook_id"] = *filter.ProviderWebhookId
		}
	}
	condition := extensions.NewSearch(nil, nil, nil, f)
	return r.BaseRepository.GetWithCondition(ctx, condition)
}
