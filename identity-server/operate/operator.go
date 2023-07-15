package operate

import (
	userStore "github.com/futugyousuzu/identity-server/user"

	mongoStore "github.com/futugyousuzu/identity-server/mongo-store"
	jwtkStore "github.com/futugyousuzu/identity-server/token"
)

type Operator struct {
	UserStore userStore.UserStore
	JwtkStore jwtkStore.JwksStore
}

func DefaultOperator() *Operator {
	o := &Operator{
		UserStore: mongoStore.NewUserStore(),
		JwtkStore: mongoStore.NewJwksStore(),
	}
	return o
}
