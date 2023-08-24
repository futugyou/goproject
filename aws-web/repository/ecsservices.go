package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IEcsServiceRepository interface {
	core.IRepository[entity.EcsServiceEntity, string]
	GetParametersByAccountId(ctx context.Context, accountId string) ([]*entity.EcsServiceEntity, error)
	GetParametersByAccountIdAndRegion(ctx context.Context, accountId string, region string) ([]*entity.EcsServiceEntity, error)
	GetParameter(ctx context.Context, accountId string, region string, key string) (*entity.EcsServiceEntity, error)
	BulkWrite(ctx context.Context, entities []entity.EcsServiceEntity) error
}
