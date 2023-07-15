package token

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JwksStore interface {
	GetPublicJwks(ctx context.Context) (string, error)
	CreateJwks(ctx context.Context, signed_key_id string) error
	GetJwkByKeyID(ctx context.Context, signed_key_id string) (jwk.Key, error)
}

type JwkModel struct {
	ID      string `bson:"_id" json:"id"`
	Payload string `bson:"payload" json:"payload"`
}
