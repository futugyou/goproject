package mongorepo

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type EcsServiceRepository struct {
	*MongoRepository[entity.EcsServiceEntity, string]
}

func NewEcsServiceRepository(config DBConfig) *EcsServiceRepository {
	baseRepo := NewMongoRepository[entity.EcsServiceEntity, string](config)
	return &EcsServiceRepository{baseRepo}
}

func (a *EcsServiceRepository) BulkWrite(ctx context.Context, entities []entity.EcsServiceEntity) error {
	models := make([]mongo.WriteModel, len(entities))
	for i := 0; i < len(entities); i++ {
		e := entities[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}

		filter := bson.D{{Key: "account_id", Value: e.AccountId}, {Key: "cluster", Value: e.Cluster}, {Key: "service_name", Value: e.ServiceName}}
		model := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpsert(true).
			SetUpdate(bson.M{
				"$set": doc,
			})
		models[i] = model
	}

	return a.BulkOperate(ctx, models)
}
