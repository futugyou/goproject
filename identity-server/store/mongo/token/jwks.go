package jwkstore

import (
	mongostore "github.com/futugyousuzu/identity-server/store/mongo"
	"github.com/futugyousuzu/identity-server/token"
)

type JwksStore struct {
	*mongostore.MongoStore
	*mongostore.InsertStore[*token.JwkModel]
	*mongostore.UpdateStore[*token.JwkModel, string]
	*mongostore.GetStore[*token.JwkModel, string]
	*mongostore.GetAllStore[*token.JwkModel]
}

func New(config mongostore.DBConfig) *JwksStore {
	baseRepo := mongostore.NewMongoStore(config)
	insertRepo := mongostore.NewInsertStore[*token.JwkModel](baseRepo)
	updateRepo := mongostore.NewUpdateStore[*token.JwkModel, string](baseRepo)
	getRepo := mongostore.NewGetStore[*token.JwkModel, string](baseRepo)
	getAllRepo := mongostore.NewGetAllStore[*token.JwkModel](baseRepo)
	return &JwksStore{baseRepo, insertRepo, updateRepo, getRepo, getAllRepo}
}
