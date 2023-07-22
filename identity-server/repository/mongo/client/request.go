package client

import (
	"github.com/futugyousuzu/identity-server/client"
	base "github.com/futugyousuzu/identity-server/repository/mongo"
)

type OAuthRequestRepository struct {
	*base.MongoRepository
	*base.InsertRepository[*client.AuthModel]
	*base.UpdateRepository[*client.AuthModel, string]
	*base.GetRepository[client.AuthModel, string]
}

func NewOAuthRequestRepository(config base.DBConfig) *OAuthRequestRepository {
	baseRepo := base.NewMongoRepository(config)
	insertRepo := base.NewInsertRepository[*client.AuthModel](baseRepo)
	updateRepo := base.NewUpdateRepository[*client.AuthModel, string](baseRepo)
	getRepo := base.NewGetRepository[client.AuthModel, string](baseRepo)
	return &OAuthRequestRepository{baseRepo, insertRepo, updateRepo, getRepo}
}
