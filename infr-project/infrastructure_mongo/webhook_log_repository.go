package infrastructure_mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func (r *WebhookLogRepository) SearchWebhookLogs(ctx context.Context, filter webhook.WebhookLogSearch) ([]webhook.WebhookLogs, error) {
	f := map[string]interface{}{}
	if len(filter.Source) > 0 {
		f["source"] = filter.Source
	}
	if len(filter.EventType) > 0 {
		f["event_type"] = filter.EventType
	}
	if len(filter.ProviderPlatformId) > 0 {
		f["provider_platform_id"] = filter.ProviderPlatformId
	}
	if len(filter.ProviderProjectId) > 0 {
		f["provider_project_id"] = filter.ProviderProjectId
	}
	if len(filter.ProviderWebhookId) > 0 {
		f["provider_webhook_id"] = filter.ProviderWebhookId
	}

	condition := extensions.NewSearch(nil, nil, nil, f)
	return r.BaseRepository.GetWithCondition(ctx, condition)
}

func (r *WebhookLogRepository) DeleteWebhookLogsByDate(ctx context.Context, filter time.Time) error {
	tenDaysAgoStr := filter.Format(time.RFC3339)
	a := new(webhook.WebhookLogs)
	c := r.Client.Database(r.DBName).Collection((*a).AggregateName())
	_, err := c.DeleteMany(context.Background(), bson.M{"happened_at": bson.M{"$lt": tenDaysAgoStr}})
	return err
}

func (r *WebhookLogRepository) InsertAndDeleteOldData(ctx context.Context, logs []webhook.WebhookLogs, filter time.Time) error {
	tenDaysAgoStr := filter.Format(time.RFC3339)
	a := new(webhook.WebhookLogs)
	c := r.Client.Database(r.DBName).Collection((*a).AggregateName())

	bulkOps := []mongo.WriteModel{
		mongo.NewDeleteManyModel().SetFilter(bson.M{"happened_at": bson.M{"$lt": tenDaysAgoStr}}),
	}

	for _, log := range logs {
		bulkOps = append(bulkOps, mongo.NewInsertOneModel().SetDocument(log))
	}

	_, err := c.BulkWrite(context.Background(), bulkOps)

	return err
}
