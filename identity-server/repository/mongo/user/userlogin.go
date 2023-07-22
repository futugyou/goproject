package userRepository

import (
	base "github.com/futugyousuzu/identity-server/repository/mongo"
	"github.com/futugyousuzu/identity-server/user"
)

type UserLoginRepository struct {
	*base.MongoRepository
	*base.InsertRepository[*user.UserLogin]
	*base.GetRepository[user.UserLogin, string]
}

func NewUserLoginRepository(config base.DBConfig) *UserLoginRepository {
	baseRepo := base.NewMongoRepository(config)
	insertRepo := base.NewInsertRepository[*user.UserLogin](baseRepo)
	getRepo := base.NewGetRepository[user.UserLogin, string](baseRepo)
	return &UserLoginRepository{baseRepo, insertRepo, getRepo}
}
