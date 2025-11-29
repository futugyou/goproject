package domain

import (
	"github.com/futugyou/domaincore/domain"
	"github.com/futugyousuzu/identity-server/pkg/security"
)

const defaultScope = "openid profile email offline_access"

type Client struct {
	domain.Aggregate
	Name         string
	RedirectUris []string
	Secret       string
	Public       bool
	Scope        string
}

type ClientOption func(*Client)

func WithScope(scope string) ClientOption {
	return func(w *Client) {
		w.Scope = scope
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

	if len(client.Scope) == 0 {
		client.Scope = defaultScope
	}

	return client, nil
}
