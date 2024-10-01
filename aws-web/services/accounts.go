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

func (a *AccountService) GetAllAccounts(ctx context.Context) []model.UserAccount {
	accounts := make([]model.UserAccount, 0)
	entities, err := a.repository.GetAll(ctx)
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

func (a *AccountService) GetAccountsByPaging(ctx context.Context, paging core.Paging) []model.UserAccount {
	accounts := make([]model.UserAccount, 0)
	entities, err := a.repository.Paging(ctx, paging)
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

func (a *AccountService) CreateAccount(ctx context.Context, account model.UserAccount) error {
	alias := account.Alias
	if len(alias) == 0 {
		log.Println("alias MUST have")
		return errors.New("alias MUST have")
	}

	accountEntity, _ := a.repository.GetAccountByAlias(ctx, alias)
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

	err := a.repository.Insert(ctx, entity)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) UpdateAccount(ctx context.Context, account model.UserAccount) error {
	entity := entity.AccountEntity{
		Id:              account.Id,
		Alias:           account.Alias,
		AccessKeyId:     account.AccessKeyId,
		SecretAccessKey: account.SecretAccessKey,
		Region:          account.Region,
		CreatedAt:       time.Now().Unix(),
	}

	err := a.repository.Update(ctx, entity, account.Id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) DeleteAccount(ctx context.Context, id string) error {
	err := a.repository.Delete(ctx, id)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (a *AccountService) GetAccountByID(ctx context.Context, id string) *model.UserAccount {
	entity, err := a.repository.Get(ctx, id)
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

func (a *AccountService) GetAccountByAlias(ctx context.Context, alias string) *model.UserAccount {
	entity, err := a.repository.GetAccountByAlias(ctx, alias)
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
