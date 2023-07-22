package client

import (
	"time"
)

type TokenModel struct {
	ID           string    `bson:"_id"`
	AccessToken  string    `bson:"access_token"`
	TokenType    string    `bson:"token_type"`
	RefreshToken string    `bson:"refresh_token"`
	Expiry       time.Time `bson:"expiry"`
}

func (TokenModel) GetType() string {
	return "oauth_tokens"
}

type AuthModel struct {
	ID                  string    `bson:"_id"`
	CodeVerifier        string    `bson:"code_verifier"`
	CodeChallenge       string    `bson:"code_challenge"`
	CodeChallengeMethod string    `bson:"code_challenge_method"`
	State               string    `bson:"state"`
	RequestURI          string    `bson:"request_uri"`
	CreateAt            time.Time `bson:"create_at"`
}

func (AuthModel) GetType() string {
	return "oauth_requests"
}
