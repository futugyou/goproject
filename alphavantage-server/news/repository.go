package news

import (
	"github.com/futugyou/alphavantage-server/core"
)

type INewsRepository interface {
	core.IRepository[NewsEntity, string]
}

type NewsRepository struct {
	*core.MongoRepository[NewsEntity, string]
}

func NewNewsRepository(config core.DBConfig) *NewsRepository {
	baseRepo := core.NewMongoRepository[NewsEntity, string](config)
	return &NewsRepository{baseRepo}
}
