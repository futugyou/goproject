package mongorepo

import (
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type AccountRepository struct {
	*MongoRepository[entity.AccountEntity, string]
}

func NewAccountRepository(config DBConfig) *AccountRepository {
	baseRepo := NewMongoRepository[entity.AccountEntity, string](config)
	return &AccountRepository{baseRepo}
}
