package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IEcsServiceRepository interface {
	core.IRepository[entity.EcsServiceEntity, string]
	BulkWrite(ctx context.Context, entities []entity.EcsServiceEntity) error
	FilterPaging(ctx context.Context, page core.Paging, filter entity.EcsServiceSearchFilter) ([]*entity.EcsServiceEntity, error)
}
