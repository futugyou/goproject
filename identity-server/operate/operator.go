package operate

import (
	base "github.com/futugyousuzu/identity-server/repository/mongo"
	jwkRepoImpl "github.com/futugyousuzu/identity-server/repository/mongo/token"
	userRepoImpl "github.com/futugyousuzu/identity-server/repository/mongo/user"

	jwksInterface "github.com/futugyousuzu/identity-server/token"
	userInterface "github.com/futugyousuzu/identity-server/user"
)

type Operator struct {
	UserRepository      userInterface.IUserRepository
	UserLoginRepository userInterface.IUserLoginRepository
	JwksRepository      jwksInterface.IJwksRepository
}

func NewOperator() *Operator {
	return &Operator{}
}

func DefaultOperator() *Operator {
	o := NewOperator()
	config := base.DBConfig{}
	jwt := jwkRepoImpl.NewJwksRepository(config)
	user := userRepoImpl.NewUserRepository(config)
	userlogin := userRepoImpl.NewUserLoginRepository(config)

	o.SetUserRepository(user)
	o.SetUserLoginRepository(userlogin)
	o.SetJwksRepository(jwt)

	return o
}

func (o *Operator) SetUserRepository(repo userInterface.IUserRepository) {
	o.UserRepository = repo
}

func (o *Operator) SetUserLoginRepository(repo userInterface.IUserLoginRepository) {
	o.UserLoginRepository = repo
}

func (o *Operator) SetJwksRepository(repo jwksInterface.IJwksRepository) {
	o.JwksRepository = repo
}
