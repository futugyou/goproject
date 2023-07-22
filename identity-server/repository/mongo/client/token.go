package client

import (
	"github.com/futugyousuzu/identity-server/client"
	base "github.com/futugyousuzu/identity-server/repository/mongo"
)

type OAuthTokenRepository struct {
	*base.MongoRepository
	*base.InsertRepository[*client.TokenModel]
	*base.UpdateRepository[*client.TokenModel, string]
	*base.GetRepository[client.TokenModel, string]
}

func NewOAuthTokenRepository(config base.DBConfig) *OAuthTokenRepository {
	baseRepo := base.NewMongoRepository(config)
	insertRepo := base.NewInsertRepository[*client.TokenModel](baseRepo)
	updateRepo := base.NewUpdateRepository[*client.TokenModel, string](baseRepo)
	getRepo := base.NewGetRepository[client.TokenModel, string](baseRepo)
	return &OAuthTokenRepository{baseRepo, insertRepo, updateRepo, getRepo}
}
