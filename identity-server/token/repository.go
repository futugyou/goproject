package token

import (
	"github.com/futugyousuzu/identity-server/core"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination mock_jwks_repo_test.go -package=token_test github.com/futugyousuzu/identity-server/token IJwksRepository

type IJwksRepository interface {
	core.IInsertRepository[*JwkModel]
	core.IGetAllRepository[JwkModel]
	core.IGetRepository[JwkModel, string]
	core.IUpdateRepository[*JwkModel, string]
}
