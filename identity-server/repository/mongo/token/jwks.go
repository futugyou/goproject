package jwksRepository

import (
	base "github.com/futugyousuzu/identity-server/repository/mongo"
	"github.com/futugyousuzu/identity-server/token"
)

type JwksRepository struct {
	*base.MongoRepository
	*base.InsertRepository[*token.JwkModel]
	*base.UpdateRepository[*token.JwkModel, string]
	*base.GetRepository[token.JwkModel, string]
	*base.GetAllRepository[token.JwkModel]
}

func NewJwksRepository(config base.DBConfig) *JwksRepository {
	baseRepo := base.NewMongoRepository(config)
	insertRepo := base.NewInsertRepository[*token.JwkModel](baseRepo)
	updateRepo := base.NewUpdateRepository[*token.JwkModel, string](baseRepo)
	getRepo := base.NewGetRepository[token.JwkModel, string](baseRepo)
	getAllRepo := base.NewGetAllRepository[token.JwkModel](baseRepo)
	return &JwksRepository{baseRepo, insertRepo, updateRepo, getRepo, getAllRepo}
}
