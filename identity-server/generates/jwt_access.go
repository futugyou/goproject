package generates

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/futugyousuzu/identity-server/operate"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// this is for go-oauth2 lib
// JWTAccessClaims jwt claims
type JWTAccessClaims struct {
	jwt.Token
}

// Valid claims verification
func (a *JWTAccessClaims) Valid() error {
	if time.Unix(a.Expiration().Unix(), 0).Before(time.Now()) {
		return errors.ErrInvalidAccessToken
	}
	return nil
}

// NewJWTAccessGenerate create to generate the jwt access token instance
func NewJWTAccessGenerate(kid string, key []byte, method jwa.SignatureAlgorithm) *JWTAccessGenerate {
	return &JWTAccessGenerate{
		SignedKeyID:  kid,
		SignedKey:    key,
		SignedMethod: method,
	}
}

// JWTAccessGenerate generate the jwt access token
type JWTAccessGenerate struct {
	SignedKeyID  string
	SignedKey    []byte
	SignedMethod jwa.SignatureAlgorithm
}

// Token based on the UUID generated token
func (a *JWTAccessGenerate) Token(ctx context.Context, data *oauth2.GenerateBasic, isGenRefresh bool) (string, string, error) {
	token_time := data.TokenInfo.GetAccessCreateAt()

	// token := jwt.New()
	// token.Set(jwt.AudienceKey, data.Client.GetID())
	// token.Set(jwt.SubjectKey, data.UserID)
	// token.Set(jwt.ExpirationKey, token_time.Add(data.TokenInfo.GetAccessExpiresIn()).Unix())
	// token.Set(jwt.IssuedAtKey, token_time)
	// token.Set(jwt.JwtIDKey, a.SignedKeyID)

	// issuer_key := os.Getenv("issuer_key")
	// if len(issuer_key) > 0 {
	// 	token.Set(jwt.IssuerKey, issuer_key)
	// }

	// if len(data.TokenInfo.GetScope()) > 0 {
	// 	token.Set("scope", data.TokenInfo.GetScope())
	// }

	jwtBuilder := jwt.NewBuilder()
	jwtBuilder.Audience([]string{data.Client.GetID()})
	jwtBuilder.Subject(data.UserID)
	jwtBuilder.Expiration(token_time.Add(data.TokenInfo.GetAccessExpiresIn()))
	jwtBuilder.IssuedAt(token_time)
	jwtBuilder.JwtID(a.SignedKeyID)

	if len(data.TokenInfo.GetScope()) > 0 {
		jwtBuilder.Claim("scope", data.TokenInfo.GetScope())
	}

	// set some claim
	operator := operate.DefaultOperator()
	userstore := operator.UserStore
	user, err := userstore.GetByUID(ctx, data.UserID)
	if err != nil {
		fmt.Printf("failed to get user data: %s\n", err)
		return "", "", err
	}

	jwtBuilder.Claim("brth", user.Birth)
	jwtBuilder.Claim("phone", user.Phone)
	jwtBuilder.Claim("name", user.Name)
	jwtBuilder.Claim("email", user.Email)

	issuer_key := os.Getenv("issuer_key")
	if len(issuer_key) > 0 {
		jwtBuilder.Issuer(issuer_key)
	}

	token, _ := jwtBuilder.Build()
	// signingKey, err := jwk.FromRaw(a.SignedKey)
	// if err != nil {
	// 	fmt.Printf("failed to create bogus JWK: %s\n", err)
	// 	return "", "", err
	// }

	// signingKey.Set(jwk.AlgorithmKey, a.SignedMethod)
	// signingKey.Set(jwk.KeyIDKey, a.SignedKey)

	store := operator.JwtkStore
	signingKey, err := store.GetJwkByKeyID(ctx, a.SignedKeyID)
	if err != nil {
		fmt.Printf("failed to create bogus JWK: %s\n", err)
		return "", "", err
	}

	// including arbitrary headers
	hdrs := jws.NewHeaders()
	hdrs.Set(`x-example`, true)

	signed, err := jwt.Sign(token, jwt.WithKey(a.SignedMethod, signingKey, jws.WithProtectedHeaders(hdrs)))

	access := string(signed)
	if err != nil {
		return "", "", err
	}
	refresh := ""

	if isGenRefresh {
		t := uuid.NewSHA1(uuid.Must(uuid.NewRandom()), []byte(access)).String()
		refresh = base64.URLEncoding.EncodeToString([]byte(t))
		refresh = strings.ToUpper(strings.TrimRight(refresh, "="))
	}

	return access, refresh, nil
}
