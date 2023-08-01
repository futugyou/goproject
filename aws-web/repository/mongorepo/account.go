package mongorepo

import (
	"context"
	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type AccountRepository struct {
	*MongoRepository[entity.AccountEntity, string]
}

func NewAccountRepository(config DBConfig) *AccountRepository {
	baseRepo := NewMongoRepository[entity.AccountEntity, string](config)
	return &AccountRepository{baseRepo}
}

func (a *AccountRepository) DeleteAll(ctx context.Context) error {
	account := new(entity.AccountEntity)
	c := a.Client.Database(a.DBName).Collection(account.GetType())
	filter := bson.D{}
	result, err := c.DeleteMany(ctx, filter)
	if err != nil {
		return err
	}

	log.Println("DeletedCount: ", result.DeletedCount)
	return nil
}
