package domain

import (
	"github.com/futugyou/domaincore/domain"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination mock_user_login_repo_test.go -package=user_test github.com/futugyousuzu/identity-server/pkg/domain  UserLoginRepository
type UserLoginRepository interface {
	domain.Repository[UserLogin]
}
