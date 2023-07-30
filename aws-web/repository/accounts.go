package repository

import (
	"github.com/futugyousuzu/goproject/awsgolang/core"
)

type IAccountRepository interface {
	core.IRepository[core.IEntity, string]
}
