package operate

import (
	base "github.com/futugyousuzu/identity-server/store/mongo"
	jwkstore "github.com/futugyousuzu/identity-server/store/mongo/token"
	userstore "github.com/futugyousuzu/identity-server/store/mongo/user"
	jwksRepository "github.com/futugyousuzu/identity-server/token"
	userRepository "github.com/futugyousuzu/identity-server/user"
)

type Operator struct {
	UserRepository      userRepository.IUserRepository
	UserLoginRepository userRepository.IUserLoginRepository
	JwksRepository      jwksRepository.IJwksRepository
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
	}

	return o
}

func (o *Operator) SetUserRepository(repo userRepository.IUserRepository) {
	o.UserRepository = repo
}

func (o *Operator) SetUserLoginRepository(repo userRepository.IUserLoginRepository) {
	o.UserLoginRepository = repo
}

func (o *Operator) SetJwksRepository(repo jwksRepository.IJwksRepository) {
	o.JwksRepository = repo
}
