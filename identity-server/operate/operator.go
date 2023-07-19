package operate

import (
	base "github.com/futugyousuzu/identity-server/store/mongo"
	jwkstore "github.com/futugyousuzu/identity-server/store/mongo/token"
	userstore "github.com/futugyousuzu/identity-server/store/mongo/user"
	jwksRepository "github.com/futugyousuzu/identity-server/token"
	userRepository "github.com/futugyousuzu/identity-server/user"

	// TODO: remove in future
	mongoStore "github.com/futugyousuzu/identity-server/mongo-store"
)

type Operator struct {
	UserRepository      userRepository.IUserRepository
	UserLoginRepository userRepository.IUserLoginRepository
	JwksRepository      jwksRepository.IJwksRepository

	// TODO: this is service, not store
	UserStore userRepository.UserStore
}

func DefaultOperator() *Operator {
	config := base.DBConfig{}
	jwt := jwkstore.NewJwksStore(config)
	user := userstore.NewUserStore(config)
	userlogin := userstore.NewUserloginStore(config)
	o := &Operator{
		UserRepository:      user,
		UserLoginRepository: userlogin,
		JwksRepository:      jwt,

		// TODO: remove in future
		UserStore: mongoStore.NewUserStore(),
	}
	return o
}
