package mongorepo

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type KeyValueRepository struct {
	*MongoRepository[entity.KeyValueEntity, string]
}

func NewKeyValueRepository(config DBConfig) *KeyValueRepository {
	baseRepo := NewMongoRepository[entity.KeyValueEntity, string](config)
	return &KeyValueRepository{baseRepo}
}

func (a *KeyValueRepository) GetValueByKey(ctx context.Context, key string) (*entity.KeyValueEntity, error) {
	keyValue := new(entity.KeyValueEntity)
	c := a.Client.Database(a.DBName).Collection(keyValue.GetType())
	filter := bson.D{{Key: "key", Value: key}}
	err := c.FindOne(ctx, filter).Decode(&keyValue)
	if err != nil {
		return nil, err
	}

	return keyValue, nil
}
