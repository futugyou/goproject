package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IParameterRepository interface {
	core.IRepository[entity.ParameterEntity, string]
	GetParametersByAccountId(ctx context.Context, accountId string) ([]*entity.ParameterEntity, error)
	GetParametersByAccountIdAndRegion(ctx context.Context, accountId string, region string) ([]*entity.ParameterEntity, error)
	GetParameter(ctx context.Context, accountId string, region string, key string) (*entity.ParameterEntity, error)
}

type IParameterLogRepository interface {
	core.IRepository[entity.ParameterLogEntity, string]
}