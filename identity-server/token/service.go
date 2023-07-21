package token

import (
	"context"

	"github.com/lestrrat-go/jwx/v2/jwk"
)

//go:generate gomockhandler -config=../gomockhandler.json  -destination ../mocks/mock_jwks_service_test.go -package=core_test github.com/futugyousuzu/identity-server/token IJwksService

type IJwksService interface {
	GetPublicJwks(ctx context.Context) (string, error)
	CreateJwks(ctx context.Context, signed_key_id string) error
	GetJwkByKeyID(ctx context.Context, signed_key_id string) (jwk.Key, error)
}
