package domain

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID        string
	Value     string
	IssuedAt  time.Time
	ExpiresAt time.Time
	ClientID  string
	UserID    string
	Scopes    []string
	Type      TokenType
	Status    TokenStatus
}

func NewIDToken(value, clientID, userID string, scopes []string, expiresAt time.Time) *Token {
	return newToken(value, expiresAt, clientID, userID, scopes, IDToken)
}

func NewAccessToken(value, clientID, userID string, scopes []string, expiresAt time.Time) *Token {
	return newToken(value, expiresAt, clientID, userID, scopes, AccessToken)
}

func NewRefreshToken(value, clientID, userID string, scopes []string, expiresAt time.Time) *Token {
	return newToken(value, expiresAt, clientID, userID, scopes, RefreshToken)
}

func newToken(value string, expiresAt time.Time, clientID string, userID string, scopes []string, tokenType TokenType) *Token {
	return &Token{
		ID:        uuid.New().String(),
		Value:     value,
		IssuedAt:  time.Now(),
		ExpiresAt: expiresAt,
		ClientID:  clientID,
		UserID:    userID,
		Scopes:    scopes,
		Type:      tokenType,
		Status:    Valid,
	}
}

type TokenType string

var AccessToken TokenType = "AccessToken"
var RefreshToken TokenType = "RefreshToken"
var IDToken TokenType = "IDToken"

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
