package domain

import (
	"github.com/futugyou/domaincore/domain"
	"github.com/google/uuid"
)

//go:generate gotests -w -all .
type User struct {
	domain.Aggregate
	Name          string
	Password      string
	Email         string
	EmailVerified bool
}

func NewUser(name, password, email string) *User {
	return &User{
		Aggregate: domain.Aggregate{
			ID: uuid.New().String(),
		},
		Name:          name,
		Password:      password,
		Email:         email,
		EmailVerified: false,
	}
}
