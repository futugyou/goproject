package mongorepo

import (
	"context"
	"log"
	"time"

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

func (a *AccountRepository) InsertMany(ctx context.Context, accounts []entity.AccountEntity) error {
	if len(accounts) == 0 {
		return nil
	}

	account := accounts[0]
	c := a.Client.Database(a.DBName).Collection(account.GetType())
	entitys := make([]interface{}, len(accounts))
	for i := 0; i < len(accounts); i++ {
		accounts[i].CreatedAt = time.Now().Unix()
		entitys[i] = accounts[i]
	}

	result, err := c.InsertMany(ctx, entitys)
	if err != nil {
		return err
	}

	log.Println("InsertedIDs: ", result.InsertedIDs)
	return nil
}
