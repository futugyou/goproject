package repository

import (
	"context"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IAccountRepository interface {
	core.IRepository[entity.AccountEntity, string]
	DeleteAll(ctx context.Context) error
	InsertMany(ctx context.Context, accounts []entity.AccountEntity) error
}
