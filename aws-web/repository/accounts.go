package repository

import (
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IAccountRepository interface {
	core.IRepository[entity.AccountEntity, string]
}
