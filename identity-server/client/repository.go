package client

import (
	"github.com/futugyousuzu/identity-server/core"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination mock_request_repo_test.go -package=client_test github.com/futugyousuzu/identity-server/client IOAuthRequestRepository
//go:generate gomockhandler -config=../gomockhandler.json  -destination mock_token_repo_test.go -package=client_test github.com/futugyousuzu/identity-server/client IOAuthTokenRepository

type IOAuthRequestRepository interface {
	core.IInsertRepository[*AuthModel]
	core.IGetRepository[AuthModel, string]
	core.IUpdateRepository[*AuthModel, string]
}

type IOAuthTokenRepository interface {
	core.IInsertRepository[*TokenModel]
	core.IGetRepository[TokenModel, string]
	core.IUpdateRepository[*TokenModel, string]
}
