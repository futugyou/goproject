package user

import (
	"context"
)

type UserStore interface {
	GetByName(ctx context.Context, name string) (User, error)
	GetByUID(ctx context.Context, uid string) (User, error)
	Login(ctx context.Context, name, password string) (UserLogin, error)
	CreateUser(ctx context.Context, user User) error
	UpdatePassword(ctx context.Context, name, password string) error
	ListUser(ctx context.Context) []User
}

type User struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	Password string `bson:"password"`
	Email    string `bson:"email"`
	Birth    string `bson:"brth"`
	Phone    string `bson:"phone"`
}

type UserLogin struct {
	ID        string `bson:"_id"`
	UserID    string `bson:"userid"`
	Timestamp int64  `bson:"timestamp"`
}
