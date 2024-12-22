package qstash

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type VerifyService service

func (s *VerifyService) Verify(ctx context.Context, request VerifyRequest) (jwt.MapClaims, error) {
	clains, err := s.verify(request.Body, request.TokenString, request.CurrentSigningKey)
	if err != nil {
		fmt.Printf("unable to verify signature with current signing key: %v", err)
		clains, err = s.verify(request.Body, request.TokenString, request.NextSigningKey)
	}

	return clains, err
}

func (s *VerifyService) verify(body []byte, tokenString, signingKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(signingKey), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if !claims.VerifyIssuer("Upstash", true) {
		return claims, fmt.Errorf("invalid issuer")
	}
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return claims, fmt.Errorf("token has expired")
	}
	if !claims.VerifyNotBefore(time.Now().Unix(), true) {
		return claims, fmt.Errorf("token is not valid yet")
	}

	bodyHash := sha256.Sum256(body)
	if claims["body"] != base64.URLEncoding.EncodeToString(bodyHash[:]) {
		return claims, fmt.Errorf("body hash does not match")
	}

	return claims, nil
}

type VerifyRequest struct {
	// Http Request Body. eg io.ReadAll(Request.Body)
	Body []byte
	// Http Request Header Upstash-Signature
	TokenString string
	// env or SigningKeysService
	CurrentSigningKey string
	// env or SigningKeysService
	NextSigningKey string
}
