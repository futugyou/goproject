package domain

import (
	"github.com/futugyou/domaincore/domain"
)

//go:generate gotests -w -all .
type UserLogin struct {
	domain.Aggregate
	UserID    string
	Timestamp int64
}
