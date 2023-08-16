package services

import (
	"context"
	"errors"
	"os"
	"time"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	"github.com/google/uuid"
)

type AccountService struct {
	repository repository.IAccountRepository
}

func NewAccountService() *AccountService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &AccountService{
		repository: mongorepo.NewAccountRepository(config),
	}
}

func (a *AccountService) GetAllAccounts() []model.UserAccount {
	accounts := make([]model.UserAccount, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return accounts
	}

	for _, entity := range entities {
		accounts = append(accounts, model.UserAccount{
			Id:              entity.Id,
			AccessKeyId:     entity.AccessKeyId,
			Alias:           entity.Alias,
			Region:          entity.Region,
			SecretAccessKey: entity.SecretAccessKey,
			CreatedAt:       time.Unix(entity.CreatedAt, 0),
		})
	}

	return accounts
}

func (a *AccountService) GetAccountsByPaging(paging core.Paging) []model.UserAccount {
	accounts := make([]model.UserAccount, 0)
	entities, err := a.repository.Paging(context.Background(), paging)
	if err != nil {
		return accounts
	}

	for _, entity := range entities {
		accounts = append(accounts, model.UserAccount{
			Id:              entity.Id,
			AccessKeyId:     entity.AccessKeyId,
			Alias:           entity.Alias,
			Region:          entity.Region,
			SecretAccessKey: entity.SecretAccessKey,
			CreatedAt:       time.Unix(entity.CreatedAt, 0),
		})
	}

	return accounts
}

func (a *AccountService) CreateAccount(account model.UserAccount) error {
	alias := account.Alias
	if len(alias) == 0 {
		log.Println("alias MUST have")
		return errors.New("alias MUST have")
	}

	accountEntity, _ := a.repository.GetAccountByAlias(context.Background(), alias)
	if accountEntity != nil && accountEntity.Alias == alias {
		log.Println("data already exist")
		return errors.New("data already exist")
	}

	entity := entity.AccountEntity{
		Id:              uuid.New().String(),
		Alias:           account.Alias,
		AccessKeyId:     account.AccessKeyId,
		SecretAccessKey: account.SecretAccessKey,
		Region:          account.Region,
		CreatedAt:       time.Now().Unix(),
	}

	err := a.repository.Insert(context.Background(), entity)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) UpdateAccount(account model.UserAccount) error {
	entity := entity.AccountEntity{
		Id:              account.Id,
		Alias:           account.Alias,
		AccessKeyId:     account.AccessKeyId,
		SecretAccessKey: account.SecretAccessKey,
		Region:          account.Region,
		CreatedAt:       time.Now().Unix(),
	}

	err := a.repository.Update(context.Background(), entity, account.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) DeleteAccount(id string) error {
	err := a.repository.Delete(context.Background(), id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) GetAccountByID(id string) *model.UserAccount {
	entity, err := a.repository.Get(context.Background(), id)
	if err != nil {
		return nil
	}

	account := &model.UserAccount{
		Id:              entity.Id,
		AccessKeyId:     entity.AccessKeyId,
		Alias:           entity.Alias,
		Region:          entity.Region,
		SecretAccessKey: entity.SecretAccessKey,
		CreatedAt:       time.Unix(entity.CreatedAt, 0),
	}

	return account
}

func (a *AccountService) GetAccountByAlias(alias string) *model.UserAccount {
	entity, err := a.repository.GetAccountByAlias(context.Background(), alias)
	if err != nil {
		return nil
	}

	account := &model.UserAccount{
		Id:              entity.Id,
		AccessKeyId:     entity.AccessKeyId,
		Alias:           entity.Alias,
		Region:          entity.Region,
		SecretAccessKey: entity.SecretAccessKey,
		CreatedAt:       time.Unix(entity.CreatedAt, 0),
	}

	return account
}
