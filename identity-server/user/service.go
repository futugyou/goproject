package user

import (
	"context"
)

type IUserService interface {
	GetByName(ctx context.Context, name string) (*User, error)
	GetByUID(ctx context.Context, uid string) (*User, error)
	Login(ctx context.Context, name, password string) (*UserLogin, error)
	CreateUser(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, name, password string) error
	ListUser(ctx context.Context) []*User
}
