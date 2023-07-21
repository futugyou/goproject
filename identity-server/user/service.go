package user

import (
	"context"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination ../mocks/mock_user_service_test.go -package=core_test github.com/futugyousuzu/identity-server/user IUserService

type IUserService interface {
	GetByName(ctx context.Context, name string) (*User, error)
	GetByUID(ctx context.Context, uid string) (*User, error)
	Login(ctx context.Context, name, password string) (*UserLogin, error)
	CreateUser(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, name, password string) error
	ListUser(ctx context.Context) []*User
}
