package domain

import (
	"context"

	"github.com/futugyou/domaincore/domain"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination mock_user_repo_test.go -package=user_test github.com/futugyousuzu/identity-server/pkg/domain  UserRepository

type UserRepository interface {
	domain.Repository[User]
	FindByName(ctx context.Context, name string) (*User, error)
}
