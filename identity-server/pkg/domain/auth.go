package domain

import (
	"slices"
	"time"
)

type Authorization struct {
	AuthorizationCode string
	UserID            string
	ClientID          string
	Scopes            []string
	ExpiresAt         time.Time
}

func NewAuthorization(authorizationCode, userID, clientID string, scopes []string, expiresAt time.Time) *Authorization {
	return &Authorization{
		AuthorizationCode: authorizationCode,
		UserID:            userID,
		ClientID:          clientID,
		Scopes:            scopes,
		ExpiresAt:         expiresAt,
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
 