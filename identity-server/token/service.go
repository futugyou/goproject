package token

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

type IJwksService interface {
	GetPublicJwks(ctx context.Context) (string, error)
	CreateJwks(ctx context.Context, signed_key_id string) error
	GetJwkByKeyID(ctx context.Context, signed_key_id string) (jwk.Key, error)
}
