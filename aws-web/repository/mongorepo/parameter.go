package mongorepo

import (
	"context"
	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
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
	cursor, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*entity.ParameterEntity, 0)
	if err = cursor.All(context.TODO(), &result); err != nil {
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
	cursor, err := c.Find(context.TODO(), filter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]*entity.ParameterEntity, 0)
	if err = cursor.All(context.TODO(), &result); err != nil {
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