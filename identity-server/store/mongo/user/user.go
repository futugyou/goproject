package userstore

import (
	"context"
	"log"

	mongostore "github.com/futugyousuzu/identity-server/store/mongo"
	"github.com/futugyousuzu/identity-server/user"
	"go.mongodb.org/mongo-driver/bson"
)

type UserStore struct {
	*mongostore.MongoStore
	*mongostore.InsertStore[*user.User]
	*mongostore.UpdateStore[*user.User, string]
	*mongostore.GetStore[*user.User, string]
	*mongostore.GetAllStore[*user.User]
}

func New(config mongostore.DBConfig) *UserStore {
	baseRepo := mongostore.NewMongoStore(config)
	insertRepo := mongostore.NewInsertStore[*user.User](baseRepo)
	updateRepo := mongostore.NewUpdateStore[*user.User, string](baseRepo)
	getRepo := mongostore.NewGetStore[*user.User, string](baseRepo)
	getAllRepo := mongostore.NewGetAllStore[*user.User](baseRepo)
	return &UserStore{baseRepo, insertRepo, updateRepo, getRepo, getAllRepo}
}

func (s *UserStore) FindByName(ctx context.Context, name string) (*user.User, error) {
	entity := new(user.User)
	c := s.Client.Database(s.DBName).Collection(entity.GetType())
	filter := bson.D{{Key: "name", Value: name}}
	err := c.FindOne(ctx, filter).Decode(&entity)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return entity, nil
}
