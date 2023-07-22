package userRepository

import (
	"context"
	"log"

	base "github.com/futugyousuzu/identity-server/repository/mongo"
	"github.com/futugyousuzu/identity-server/user"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	*base.MongoRepository
	*base.InsertRepository[*user.User]
	*base.UpdateRepository[*user.User, string]
	*base.GetRepository[user.User, string]
	*base.GetAllRepository[user.User]
}

func NewUserRepository(config base.DBConfig) *UserRepository {
	baseRepo := base.NewMongoRepository(config)
	insertRepo := base.NewInsertRepository[*user.User](baseRepo)
	updateRepo := base.NewUpdateRepository[*user.User, string](baseRepo)
	getRepo := base.NewGetRepository[user.User, string](baseRepo)
	getAllRepo := base.NewGetAllRepository[user.User](baseRepo)
	return &UserRepository{baseRepo, insertRepo, updateRepo, getRepo, getAllRepo}
}

func (s *UserRepository) FindByName(ctx context.Context, name string) (*user.User, error) {
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
