package expected

import (
	"github.com/futugyou/alphavantage-server/core"
)

type IExpectedRepository interface {
	core.IRepository[ExpectedEntity, string]
}

type ExpectedRepository struct {
	*core.MongoRepository[ExpectedEntity, string]
}

func NewExpectedRepository(config core.DBConfig) *ExpectedRepository {
	baseRepo := core.NewMongoRepository[ExpectedEntity, string](config)
	return &ExpectedRepository{baseRepo}
}
