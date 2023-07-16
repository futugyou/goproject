package token

import (
	"github.com/futugyousuzu/identity-server/core"
)

type IJwksRepository interface {
	core.IInsertRepository[*JwkModel]
	core.IGetAllRepository[*JwkModel]
	core.IGetRepository[*JwkModel, string]
	core.IUpdateRepository[*JwkModel, string]
}
