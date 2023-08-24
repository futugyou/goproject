package services

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	"golang.org/x/exp/slices"
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

func (e *EcsClusterService) GetAllServices(paging core.Paging, filter model.EcsClusterFilter) ([]model.EcsClusterViewModel, error) {
	accounts := make([]model.UserAccount, 0)
	entityfilter := entity.EcsServiceSearchFilter{}
	accountService := NewAccountService()
	if len(filter.AccountId) > 0 {
		entityfilter.AccountId = filter.AccountId
		account := accountService.GetAccountByID(filter.AccountId)
		if account == nil {
			return nil, errors.New("account not found")
		}
		accounts = append(accounts, *account)
	} else {
		accounts = accountService.GetAllAccounts()
		if len(accounts) == 0 {
			return nil, errors.New("account not found")
		}
	}

	entities, err := e.repository.FilterPaging(context.Background(), paging, entityfilter)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result := make([]model.EcsClusterViewModel, 0)
	for _, entity := range entities {
		idx := slices.IndexFunc(accounts, func(c model.UserAccount) bool { return c.Id == entity.AccountId })
		alias := ""
		if idx != -1 && idx < len(accounts) {
			alias = accounts[idx].Alias
		}
		e := model.EcsClusterViewModel{
			ClusterName:  entity.Cluster,
			ClusterArn:   entity.ClusterArn,
			Service:      entity.ServiceName,
			ServiceArn:   entity.ServiceNameArn,
			RoleArn:      entity.RoleArn,
			AccountAlias: alias,
			OperateAt:    entity.OperateAt,
		}

		result = append(result, e)
	}
	return result, nil
}
