package domain

import (
	"slices"
	"time"

	"github.com/futugyou/domaincore/domain"
	"github.com/futugyousuzu/identity-server/pkg/security"
	"github.com/google/uuid"
)

type Authorization struct {
	domain.Aggregate
	UserID    string
	ClientID  string
	Scopes    []string
	ExpiresAt time.Time
}

func NewAuthorization(userID, clientID string, scopes []string, expiresAt time.Time) *Authorization {
	code, err := security.GenerateRandomString(24)
	if err != nil {
		code = uuid.New().String()
	}

	return &Authorization{
		Aggregate: domain.Aggregate{
			ID: code,
		},
		UserID:    userID,
		ClientID:  clientID,
		Scopes:    scopes,
		ExpiresAt: expiresAt,
	}
}

func (a *Authorization) ValidateExpired() bool {
	return a.ExpiresAt.After(time.Now())
}

func (a *Authorization) ValidateScope(scopes []string) bool {
	for _, scope := range scopes {
		if !slices.Contains(a.Scopes, scope) {
			return false
		}
	}

	return true
}

func (a *Authorization) RevokeAuthorization() {
	a.ExpiresAt = time.Now()
}
