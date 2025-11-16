package infrastructure

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	domaincore "github.com/futugyou/domaincore/domain"

	"github.com/futugyou/domaincore/mongoimpl"
	"github.com/futugyou/platformservice/domain"
)

type WebhookLogRepository struct {
	mongoimpl.BaseCRUD[domain.WebhookLogs]
}

func NewWebhookLogRepository(client *mongo.Client, config mongoimpl.DBConfig) *WebhookLogRepository {
	if config.CollectionName == "" {
		config.CollectionName = "platform_webhook_logs"
	}

	getID := func(a domain.WebhookLogs) string { return a.AggregateId() }

	return &WebhookLogRepository{
		BaseCRUD: *mongoimpl.NewBaseCRUD(client, config, getID),
	}
}

func (r *WebhookLogRepository) SearchWebhookLogs(ctx context.Context, filter domain.WebhookLogSearch) ([]domain.WebhookLogs, error) {
	f := map[string]any{}
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

	condition := domaincore.NewQueryOptions(nil, nil, nil, f)
	return r.BaseCRUD.Find(ctx, condition)
}

func (r *WebhookLogRepository) DeleteWebhookLogsByDate(ctx context.Context, filter time.Time) error {
	tenDaysAgoStr := filter.Format(time.RFC3339)
	c := r.Client.Database(r.DBName).Collection(r.CollectionName)
	_, err := c.DeleteMany(context.Background(), bson.M{"happened_at": bson.M{"$lt": tenDaysAgoStr}})
	return err
}

func (r *WebhookLogRepository) InsertAndDeleteOldData(ctx context.Context, logs []domain.WebhookLogs, filter time.Time) error {
	tenDaysAgoStr := filter.Format(time.RFC3339)
	c := r.Client.Database(r.DBName).Collection(r.CollectionName)

	bulkOps := []mongo.WriteModel{
		mongo.NewDeleteManyModel().SetFilter(bson.M{"happened_at": bson.M{"$lt": tenDaysAgoStr}}),
	}

	for _, log := range logs {
		bulkOps = append(bulkOps, mongo.NewInsertOneModel().SetDocument(log))
	}

	_, err := c.BulkWrite(context.Background(), bulkOps)

	return err
}
