package user

import (
	"github.com/futugyou/extensions"

	"github.com/futugyou/domaincore/domain"

	"github.com/google/uuid"
)

var defaultScopes []string = []string{"openid", "profile", "email", "offline_access"}

//go:generate gotests -w -all .
type User struct {
	domain.Aggregate
	Name          string
	Password      string
	Email         string
	EmailVerified bool
	Roles         []string
	Scopes        []string
}

func NewUser(name, password, email string, scopes []string) *User {
	user := &User{
		Aggregate: domain.Aggregate{
			ID: uuid.New().String(),
		},
		Name:          name,
		Password:      password,
		Email:         email,
		EmailVerified: false,
		Roles:         []string{},
		Scopes:        scopes,
	}

	if len(user.Scopes) == 0 {
		user.Scopes = defaultScopes
	}

	return user
}

func (u *User) GrantRole(roles []string) {
	u.Roles = extensions.MergeDeduplication(u.Roles, roles)
}

func (u *User) RevokeRole(role string) {
	if role == "" {
		return
	}

	for i, s := range u.Roles {
		if s == role {
			copy(u.Roles[i:], u.Roles[i+1:])
			u.Roles = u.Roles[:len(u.Roles)-1]
			return
		}
	}
}

func (u *User) GrantAuthorization(scopes []string) {
	u.Scopes = extensions.MergeDeduplication(u.Scopes, scopes)
}

func (u *User) RevokeAuthorization(scope string) {
	if scope == "" {
		return
	}

	for i, s := range u.Scopes {
		if s == scope {
			copy(u.Scopes[i:], u.Scopes[i+1:])
			u.Scopes = u.Scopes[:len(u.Scopes)-1]
			return
		}
	}
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func (u *User) HasScope(scope string) bool {
	for _, s := range u.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}
