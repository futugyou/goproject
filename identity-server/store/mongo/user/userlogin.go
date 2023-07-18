package userstore

import (
	mongostore "github.com/futugyousuzu/identity-server/store/mongo"
	"github.com/futugyousuzu/identity-server/user"
)

type UserloginStore struct {
	*mongostore.MongoStore
	*mongostore.InsertStore[*user.UserLogin]
	*mongostore.GetStore[*user.UserLogin, string]
}

func NewUserloginStore(config mongostore.DBConfig) *UserloginStore {
	baseRepo := mongostore.NewMongoStore(config)
	insertRepo := mongostore.NewInsertStore[*user.UserLogin](baseRepo)
	getRepo := mongostore.NewGetStore[*user.UserLogin, string](baseRepo)
	return &UserloginStore{baseRepo, insertRepo, getRepo}
}
