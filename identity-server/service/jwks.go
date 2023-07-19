package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/futugyousuzu/identity-server/operate"
	"github.com/futugyousuzu/identity-server/token"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
)

type JwksService struct {
	token.IJwksRepository
}

func NewJwksService(operator *operate.Operator) *JwksService {
	return &JwksService{operator.JwksRepository}
}

func (u *JwksService) GetPublicJwks(ctx context.Context) (string, error) {
	jwtRepo := u.IJwksRepository
	models, err := jwtRepo.GetAll(ctx)
	if err != nil {
		return "", err
	}

	s := make([]string, len(models))
	privset := jwk.NewSet()
	for i, v := range models {
		s[i] = v.Payload
		key, err := jwk.ParseKey([]byte(v.Payload))
		if err != nil {
			panic(err)
		}

		privset.AddKey(key)
	}

	v, err := jwk.PublicSetOf(privset)
	if err != nil {
		panic(err)
	}

	buf, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal key into JSON: %s", err)
	}

	return string(buf), err
}

func (u *JwksService) CreateJwks(ctx context.Context, signed_key_id string) error {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate new rsa private key: %s", err)
	}

	key, err := jwk.FromRaw(raw)
	if err != nil {
		return fmt.Errorf("failed to create RSA key: %s", err)
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		return fmt.Errorf("expected jwk.RSAPrivateKey, got %T", err)
	}

	key.Set(jwk.KeyIDKey, signed_key_id)
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	key.Set(`my-custom-field`, `unbelievable-value`)

	buf, err := json.MarshalIndent(key, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal key into JSON: %s", err)
	}

	jwkModel := token.JwkModel{
		ID:      signed_key_id,
		Payload: string(buf),
	}

	jwtRepo := u.IJwksRepository
	return jwtRepo.Update(ctx, &jwkModel, signed_key_id)
}

func (u *JwksService) GetJwkByKeyID(ctx context.Context, signed_key_id string) (jwk.Key, error) {
	jwtRepo := u.IJwksRepository
	jwkModel, err := jwtRepo.Get(ctx, signed_key_id)

	if err != nil {
		return nil, err
	}

	return jwk.ParseKey([]byte(jwkModel.Payload))
}
