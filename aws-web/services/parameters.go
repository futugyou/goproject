package services

import (
	"context"
	"os"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
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
	parameters := make([]model.ParameterViewModel, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		parameters = append(parameters, model.ParameterViewModel{
			Id:        entity.Id,
			AccountId: entity.AccountId,
			Region:    entity.Region,
			Key:       entity.Key,
			// Value:     entity.Value, // list not need value
			Version:   entity.Version,
			OperateAt: time.Unix(entity.OperateAt, 0),
		})
	}

	return parameters
}

func (a *ParameterService) GetParametersByPaging(paging core.Paging) []model.ParameterViewModel {
	parameters := make([]model.ParameterViewModel, 0)
	entities, err := a.repository.Paging(context.Background(), paging)
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		parameters = append(parameters, model.ParameterViewModel{
			Id:        entity.Id,
			AccountId: entity.AccountId,
			Region:    entity.Region,
			Key:       entity.Key,
			// Value:     entity.Value, // list not need value
			Version:   entity.Version,
			OperateAt: time.Unix(entity.OperateAt, 0),
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
