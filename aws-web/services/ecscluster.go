package services

import (
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type EcsClusterService struct {
	repository repository.IEcsServiceRepository
}

func NewEcsClusterService() *EcsClusterService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &EcsClusterService{
		repository: mongorepo.NewEcsServiceRepository(config),
	}
}

func (e *EcsClusterService) GetAllServices(filter *model.EcsClusterFilter) ([]model.EcsClusterViewModel, error) {
	var account *model.UserAccount
	accountService := NewAccountService()
	if filter != nil && len(filter.AccountId) > 0 {
		account = accountService.GetAccountByID(filter.AccountId)
		if account == nil {
			return nil, errors.New("account not found")
		}
	} else {
		accounts := accountService.GetAllAccounts()
		if len(accounts) == 0 {
			return nil, errors.New("account not found")
		}

		account = &accounts[0]
	}

	awsenv.CfgForVercelWithRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	ecs.NewFromConfig(awsenv.Cfg)
	return nil, nil
}
