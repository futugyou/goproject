package repository

import (
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
)

type IKeyValueRepository interface {
	core.IRepository[entity.KeyValueEntity, string]
}
