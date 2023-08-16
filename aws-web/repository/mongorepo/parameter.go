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

type ParameterRepository struct {
	*MongoRepository[entity.ParameterEntity, string]
}

func NewParameterRepository(config DBConfig) *ParameterRepository {
	baseRepo := NewMongoRepository[entity.ParameterEntity, string](config)
	return &ParameterRepository{baseRepo}
}

func (a *ParameterRepository) GetParametersByAccountId(ctx context.Context, accountId string) ([]*entity.ParameterEntity, error) {
	parameter := new(entity.ParameterEntity)
	c := a.Client.Database(a.DBName).Collection(parameter.GetType())
	filter := bson.D{{Key: "account_id", Value: accountId}}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*entity.ParameterEntity, 0)
	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}

func (a *ParameterRepository) GetParametersByAccountIdAndRegion(ctx context.Context, accountId string, region string) ([]*entity.ParameterEntity, error) {
	parameter := new(entity.ParameterEntity)
	c := a.Client.Database(a.DBName).Collection(parameter.GetType())
	filter := bson.D{{Key: "account_id", Value: accountId}, {Key: "region", Value: region}}
	cursor, err := c.Find(ctx, filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*entity.ParameterEntity, 0)
	if err = cursor.All(ctx, &result); err != nil {
		log.Println(err)
		return nil, err
	}

	for _, data := range result {
		cursor.Decode(&data)
	}

	return result, nil
}

func (a *ParameterRepository) GetParameter(ctx context.Context, accountId string, region string, key string) (*entity.ParameterEntity, error) {
	entity := new(entity.ParameterEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())

	filter := bson.D{{Key: "account_id", Value: accountId}, {Key: "region", Value: region}, {Key: "key", Value: key}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}

type ParameterLogRepository struct {
	*MongoRepository[entity.ParameterLogEntity, string]
}

func NewParameterLogRepository(config DBConfig) *ParameterLogRepository {
	baseRepo := NewMongoRepository[entity.ParameterLogEntity, string](config)
	return &ParameterLogRepository{baseRepo}
}

func (a *ParameterRepository) BulkWrite(ctx context.Context, entities []entity.ParameterEntity) error {
	models := make([]mongo.WriteModel, len(entities))
	for i := 0; i < len(entities); i++ {
		e := entities[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}

		filter := bson.D{{Key: "account_id", Value: e.AccountId}, {Key: "region", Value: e.Region}, {Key: "key", Value: e.Key}}
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

func (a *ParameterLogRepository) BulkWrite(ctx context.Context, entities []entity.ParameterLogEntity) error {
	models := make([]mongo.WriteModel, len(entities))
	for i := 0; i < len(entities); i++ {
		e := entities[i]
		doc, err := flatbson.Flatten(e)
		if err != nil {
			log.Println("BulkWrite: ", i, err)
			continue
		}

		filter := bson.D{{Key: "account_id", Value: e.AccountId}, {Key: "region", Value: e.Region}, {Key: "key", Value: e.Key}, {Key: "version", Value: e.Version}}
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

func (a *ParameterRepository) FilterPaging(ctx context.Context, page core.Paging, filter entity.ParameterSearchFilter) ([]*entity.ParameterEntity, error) {
	result := make([]*entity.ParameterEntity, 0)
	entity := new(entity.ParameterEntity)
	c := a.Client.Database(a.DBName).Collection((*entity).GetType())

	filters := make([]bson.M, 0)

	if len(filter.AccountId) > 0 {
		filters = append(filters, bson.M{"account_id": filter.AccountId})
	}

	if len(filter.Key) > 0 {
		filters = append(filters, bson.M{"key": bson.M{"$regex": filter.Key, "$options": "im"}})
	}

	if len(filter.Region) > 0 {
		filters = append(filters, bson.M{"region": filter.Region})
	}

	bsonfilter := bson.M{}
	if len(filters) > 1 {
		bsonfilter = bson.M{"$and": filters}
	}

	var skip int64 = (page.Page - 1) * page.Limit
	op := options.Find().SetLimit(page.Limit).SetSkip(skip)
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
