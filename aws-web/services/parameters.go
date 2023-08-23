package services

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
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
	var searchFilter entity.ParameterSearchFilter = entity.ParameterSearchFilter{
		Key:    filter.Key,
		Region: filter.Region,
	}

	if account != nil {
		searchFilter.AccountId = account.Id
		accounts = append(accounts, *account)
	}

	entities, err := a.repository.FilterPaging(context.Background(), paging, searchFilter)
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

func (a *ParameterService) GetParameterByID(id string) *model.ParameterDetailViewModel {
	// get parameter from db
	ctx := context.Background()
	entity, err := a.repository.GetByObjectId(ctx, id)
	if err != nil {
		return nil
	}

	parameter := &model.ParameterDetailViewModel{
		Id:        entity.Id,
		AccountId: entity.AccountId,
		Region:    entity.Region,
		Key:       entity.Key,
		Value:     entity.Value,
		Version:   entity.Version,
		OperateAt: time.Unix(entity.OperateAt, 0),
	}

	// fill account alias
	accountService := NewAccountService()
	account := accountService.GetAccountByID(entity.AccountId)
	if account == nil {
		return parameter
	}

	parameter.AccountAlias = account.Alias

	// get parameter log from db
	logs, err := a.logRepository.GetParameterLogs(ctx, entity.AccountId, entity.Region, entity.Key)
	if err == nil && len(logs) > 0 {
		history := make([]model.ParameterViewModel, len(logs))
		for i := 0; i < len(logs); i++ {
			history[i] = model.ParameterViewModel{
				Id:           logs[i].Id,
				AccountId:    logs[i].AccountId,
				AccountAlias: account.Alias,
				Region:       logs[i].Region,
				Key:          logs[i].Key,
				Value:        logs[i].Value,
				Version:      logs[i].Version,
				OperateAt:    time.Unix(logs[i].OperateAt, 0),
			}
		}

		parameter.History = history
	}

	// get parameter from aws
	awsenv.CfgForVercelWithRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	details, err := a.getParametersDatail([]string{entity.Key})
	if err != nil || len(details) == 0 {
		return parameter
	}

	current := &model.ParameterViewModel{
		Key:       *details[0].Name,
		Value:     *details[0].Value,
		Version:   strconv.FormatInt(details[0].Version, 10),
		OperateAt: *details[0].LastModifiedDate,
	}

	parameter.Current = current

	return parameter
}

func (a *ParameterService) CompareParameterByIDs(sourceid string, destid string) []model.CompareViewModel {
	ctx := context.Background()
	source, err := a.repository.GetByObjectId(ctx, sourceid)
	if err != nil {
		return nil
	}

	dest, err := a.repository.GetByObjectId(ctx, destid)
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

func (a *ParameterService) SyncParameterByID(id string) error {
	// get parameter from db
	ctx := context.Background()
	parameter, err := a.repository.GetByObjectId(ctx, id)
	if err != nil {
		return err
	}

	// account
	accountService := NewAccountService()
	account := accountService.GetAccountByID(parameter.AccountId)
	if account == nil {
		return errors.New("account not found")
	}

	// get parameter from aws
	awsenv.CfgForVercelWithRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	details, err := a.getParametersDatail([]string{parameter.Key})
	if err != nil || len(details) == 0 {
		return errors.New("ssm not found")
	}

	currVersion, _ := strconv.ParseInt(parameter.Version, 10, 64)
	if details[0].Version <= currVersion {
		return nil
	}

	// update parameter
	modified := details[0].LastModifiedDate
	if modified != nil {
		parameter.OperateAt = modified.Unix()
	}

	parameter.Version = strconv.FormatInt(details[0].Version, 10)
	parameter.Value = *details[0].Value
	err = a.repository.Update(ctx, *parameter, parameter.Id)
	if err != nil {
		return err
	}

	// update parameter log
	log := entity.ParameterLogEntity{
		AccountId: parameter.AccountId,
		Region:    parameter.Region,
		Key:       parameter.Key,
		Value:     parameter.Value,
		Version:   parameter.Version,
		OperateAt: parameter.OperateAt,
	}

	return a.logRepository.Insert(ctx, log)
}
