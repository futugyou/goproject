package domain

import (
	"time"

	"github.com/futugyou/domaincore/domain"
	"github.com/google/uuid"
)

type Token struct {
	domain.Aggregate
	AuthorizationCode     string
	ClientID              string
	UserID                string
	IssuedAt              time.Time
	ExpiresAt             time.Time
	Scopes                []string
	AccessToken           string
	IDToken               string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
	Status                TokenStatus
}

func NewToken(code, clientID, userID string, scopes []string, token string, expiresAt time.Time) *Token {
	return &Token{
		Aggregate: domain.Aggregate{
			ID: uuid.New().String(),
		},
		AuthorizationCode: code,
		ClientID:          clientID,
		UserID:            userID,
		IssuedAt:          time.Now(),
		ExpiresAt:         expiresAt,
		Scopes:            scopes,
		Status:            Valid,
	}
}

type TokenStatus string

var Valid TokenStatus = "Valid"
var Revoked TokenStatus = "Revoked"
var Expired TokenStatus = "Expired"

func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

func (t *Token) Invalidate() {
	t.Status = Expired
	t.ExpiresAt = time.Now()
}
