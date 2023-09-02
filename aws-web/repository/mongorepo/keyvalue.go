package mongorepo

import (
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type KeyValueRepository struct {
	*MongoRepository[entity.KeyValueEntity, string]
}

func NewKeyValueRepository(config DBConfig) *KeyValueRepository {
	baseRepo := NewMongoRepository[entity.KeyValueEntity, string](config)
	return &KeyValueRepository{baseRepo}
}
