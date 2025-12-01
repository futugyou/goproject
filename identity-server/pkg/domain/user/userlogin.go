package domain

import (
	"time"

	"github.com/futugyou/domaincore/domain"
	"github.com/google/uuid"
)

//go:generate gotests -w -all .
type UserLogin struct {
	domain.Aggregate
	UserID    string
	Timestamp int64
}

func NewUserLogin(uid string) *UserLogin {
	return &UserLogin{
		Aggregate: domain.Aggregate{
			ID: uuid.New().String(),
		},
		UserID:    uid,
		Timestamp: time.Now().Unix(),
	}
}
