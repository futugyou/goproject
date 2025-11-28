package domain

import (
	"github.com/futugyou/domaincore/domain"
)

//go:generate gotests -w -all .
type User struct {
	domain.Aggregate
	Name     string
	Password string
	Email    string
}
