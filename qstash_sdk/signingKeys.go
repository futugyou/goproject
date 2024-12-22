package qstash

import (
	"context"
)

type SigningKeysService service

func (s *SigningKeysService) GetSigningKeys(ctx context.Context) (*SigningKeys, error) {
	path := "/v2/keys"
	result := &SigningKeys{}
	if err := s.client.http.Get(ctx, path, result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *SigningKeysService) RotateSigningKeys(ctx context.Context) (*SigningKeys, error) {
	path := "/v2/keys/rotate"
	result := &SigningKeys{}
	if err := s.client.http.Post(ctx, path, nil, result); err != nil {
		return nil, err
	}

	return result, nil
}

type SigningKeys struct {
	Current string `json:"current"`
	Next    string `json:"next"`
}
