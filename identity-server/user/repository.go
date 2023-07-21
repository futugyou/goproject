package user

import (
	"context"

	"github.com/futugyousuzu/identity-server/core"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination ../mocks/mock_user_repo_test.go -package=core_test github.com/futugyousuzu/identity-server/user IUserRepository

type IUserRepository interface {
	core.IInsertRepository[*User]
	core.IGetAllRepository[User]
	core.IGetRepository[User, string]
	core.IUpdateRepository[*User, string]
	FindByName(ctx context.Context, name string) (*User, error)
}

//go:generate gomockhandler -config=../gomockhandler.json  -destination ../mocks/mock_user_login_repo_test.go -package=core_test github.com/futugyousuzu/identity-server/user IUserLoginRepository

type IUserLoginRepository interface {
	core.IInsertRepository[*UserLogin]
	core.IGetRepository[UserLogin, string]
}
