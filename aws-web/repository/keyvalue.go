package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IKeyValueRepository interface {
	core.IRepository[entity.KeyValueEntity, string]
	GetValueByKey(ctx context.Context, key string) (*entity.KeyValueEntity, error)
}
