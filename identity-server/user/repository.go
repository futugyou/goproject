package user

import (
	"context"

	"github.com/futugyousuzu/identity-server/core"
)

type IUserRepository interface {
	core.IInsertRepository[*User]
	core.IGetAllRepository[*User]
	core.IGetRepository[*User, string]
	core.IUpdateRepository[*User, string]
	FindByName(ctx context.Context, name string) (*User, error)
}

type IUserLoginRepository interface {
	core.IInsertRepository[*UserLogin]
	core.IGetRepository[*UserLogin, string]
}
