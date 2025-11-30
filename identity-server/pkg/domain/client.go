package domain

import (
	"github.com/futugyou/domaincore/domain"
	"github.com/futugyousuzu/identity-server/pkg/security"
)

var defaultScopes []string = []string{"openid", "profile", "email", "offline_access"}

type Client struct {
	domain.Aggregate
	Name         string
	RedirectUris []string
	Secret       string
	Public       bool
	Scopes       []string
}

type ClientOption func(*Client)

func WithScopes(scopes []string) ClientOption {
	return func(w *Client) {
		w.Scopes = mergeDeduplication(w.Scopes, scopes)
	}
}

func WithPublic(public bool) ClientOption {
	return func(w *Client) {
		w.Public = public
	}
}

func NewClient(name string, redirectUris []string, opts ...ClientOption) (*Client, error) {
	id, secret, err := security.GenerateClientCredentials()
	if err != nil {
		return nil, err
	}
	client := &Client{
		Aggregate: domain.Aggregate{
			ID: id,
		},
		Name:         name,
		RedirectUris: redirectUris,
		Secret:       secret,
	}

	for _, opt := range opts {
		opt(client)
	}

	if len(client.Scopes) == 0 {
		client.Scopes = defaultScopes
	}

	return client, nil
}
