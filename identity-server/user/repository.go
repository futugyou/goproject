package user

import (
	"github.com/futugyousuzu/identity-server/core"
)

type IUserRepository interface {
	core.IInsertRepository[*User]
	core.IGetAllRepository[*User]
	core.IGetRepository[*User, string]
	core.IUpdateRepository[*User, string]
}

type IUserLoginRepository interface {
	core.IInsertRepository[*UserLogin]
	core.IGetRepository[*UserLogin, string]
}
