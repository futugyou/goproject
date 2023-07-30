package services

import (
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
)

type UserAccount struct {
	Id              string `json:"id"`
	Alias           string `json:"alias"`
	AccessKeyId     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	CreatedAt       int    `json:"createdAt"`
}

type AccountService struct {
	repository repository.IAccountRepository
}

func NewAccountService() *AccountService {
	// TODO: db config
	config := mongorepo.DBConfig{}
	return &AccountService{
		repository: mongorepo.NewAccountRepository(config),
	}
}

func (a *AccountService) GetAllAccounts() []UserAccount {
	accounts := make([]UserAccount, 0)
	return accounts
}
