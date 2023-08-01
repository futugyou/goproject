package services

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"log"

	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
)

type UserAccount struct {
	Id              string    `json:"id"`
	Alias           string    `json:"alias"`
	AccessKeyId     string    `json:"accessKeyId"`
	SecretAccessKey string    `json:"secretAccessKey"`
	Region          string    `json:"region"`
	CreatedAt       time.Time `json:"createdAt"`
}

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

func (a *AccountService) GetAllAccounts() []UserAccount {
	accounts := make([]UserAccount, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return accounts
	}

	for _, entity := range entities {
		accounts = append(accounts, UserAccount{
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

func (a *AccountService) AccountInit() {
	result := make([]entity.AccountEntity, 0)
	var accounts []byte
	var err error

	if accounts, err = os.ReadFile("./data/accounts.json"); err != nil {
		log.Println(err)
		return
	}

	if err = json.Unmarshal(accounts, &result); err != nil {
		log.Println(err)
		return
	}

	if err = a.repository.DeleteAll(context.Background()); err != nil {
		log.Println(err)
		return
	}

	for _, acc := range result {
		acc.CreatedAt = time.Now().Unix()
	}

	if err = a.repository.InsertMany(context.Background(), result); err != nil {
		log.Println(err)
	}
}
