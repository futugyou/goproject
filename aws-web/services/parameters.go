package services

import (
	"context"
	"os"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
	"golang.org/x/exp/slices"
)

type ParameterService struct {
	repository    repository.IParameterRepository
	logRepository repository.IParameterLogRepository
}

func NewParameterService() *ParameterService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &ParameterService{
		repository:    mongorepo.NewParameterRepository(config),
		logRepository: mongorepo.NewParameterLogRepository(config),
	}
}

func (a *ParameterService) GetAllParameters() []model.ParameterViewModel {
	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts()

	parameters := make([]model.ParameterViewModel, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		idx := slices.IndexFunc(accounts, func(c model.UserAccount) bool { return c.Id == entity.AccountId })
		parameters = append(parameters, model.ParameterViewModel{
			Id:           entity.Id,
			AccountId:    entity.AccountId,
			AccountAlias: accounts[idx].Alias,
			Region:       entity.Region,
			Key:          entity.Key,
			// Value:     entity.Value, // list not need value
			Version:   entity.Version,
			OperateAt: time.Unix(entity.OperateAt, 0),
		})
	}

	return parameters
}

func (a *ParameterService) GetParametersByCondition(paging core.Paging, filter model.ParameterFilter) []model.ParameterViewModel {
	accountService := NewAccountService()
	var account *model.UserAccount
	var accounts = make([]model.UserAccount, 0)
	if len(filter.AccountAlias) > 0 {
		account = accountService.GetAccountByAlias(filter.AccountAlias)
	} else {
		accounts = accountService.GetAllAccounts()
	}

	parameters := make([]model.ParameterViewModel, 0)
	var SearchFilter entity.ParameterSearchFilter = entity.ParameterSearchFilter{
		Key:    filter.Key,
		Region: filter.Region,
	}

	if account != nil {
		SearchFilter.AccountId = account.Id
		accounts = append(accounts, *account)
	}

	entities, err := a.repository.FilterPaging(context.Background(), paging, SearchFilter)
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		idx := slices.IndexFunc(accounts, func(c model.UserAccount) bool { return c.Id == entity.AccountId })
		alias := ""
		if idx != -1 && idx < len(accounts) {
			alias = accounts[idx].Alias
		}

		parameters = append(parameters, model.ParameterViewModel{
			Id:           entity.Id,
			AccountId:    entity.AccountId,
			AccountAlias: alias,
			Region:       entity.Region,
			Key:          entity.Key,
			Version:      entity.Version,
			OperateAt:    time.Unix(entity.OperateAt, 0),
		})
	}

	return parameters
}

func (a *ParameterService) GetParameterByID(id string) *model.ParameterViewModel {
	entity, err := a.repository.Get(context.Background(), id)
	if err != nil {
		return nil
	}

	parameter := &model.ParameterViewModel{
		Id:        entity.Id,
		AccountId: entity.AccountId,
		Region:    entity.Region,
		Key:       entity.Key,
		Value:     entity.Value,
		Version:   entity.Version,
		OperateAt: time.Unix(entity.OperateAt, 0),
	}

	return parameter
}

func (a *ParameterService) CompareParameterByIDs(sourceid string, destid string) []model.CompareViewModel {
	source, err := a.repository.Get(context.Background(), sourceid)
	if err != nil {
		return nil
	}

	dest, err := a.repository.Get(context.Background(), destid)
	if err != nil {
		return nil
	}

	result := []model.CompareViewModel{{
		Key:     source.Key,
		Value:   source.Value,
		Version: source.Version,
	}, {
		Key:     dest.Key,
		Value:   dest.Value,
		Version: dest.Version,
	}}

	return result
}
