package mongorepo

import (
	"context"
	"log"

	"github.com/chidiwilliams/flatbson"
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		filter := bson.D{{Key: "account_id", Value: e.AccountId}, {Key: "cluster_arn", Value: e.ClusterArn}, {Key: "service_name_arn", Value: e.ServiceNameArn}}
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
func (a *EcsServiceRepository) FilterPaging(ctx context.Context, page core.Paging, filter entity.EcsServiceSearchFilter) ([]*entity.EcsServiceEntity, error) {
	result := make([]*entity.EcsServiceEntity, 0)
	entity := new(entity.EcsServiceEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())

	filters := make([]bson.M, 0)

	if len(filter.AccountId) > 0 {
		filters = append(filters, bson.M{"account_id": filter.AccountId})
	}

	bsonfilter := bson.M{}
	if len(filters) == 1 {
		bsonfilter = filters[0]
	} else if len(filters) > 1 {
		bsonfilter = bson.M{"$and": filters}
	}

	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip).SetSort(bson.D{{Key: "operate_at", Value: 1}})
	cursor, err := c.Find(ctx, bsonfilter, op)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}
