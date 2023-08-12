package services

import (
	"context"
	"os"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
)

type ParameterViewModel struct {
	Id        string    `bson:"_id,omitempty"`
	AccountId string    `bson:"account_id"`
	Region    string    `bson:"region"`
	Key       string    `bson:"key"`
	Value     string    `bson:"value"`
	Version   string    `bson:"version"`
	OperateAt time.Time `bson:"operate_at"`
}

type ParameterService struct {
	repository repository.IParameterRepository
}

func NewParameterService() *ParameterService {
	config := mongorepo.DBConfig{
		DBName:        os.Getenv("db_name"),
		ConnectString: os.Getenv("mongodb_url"),
	}

	return &ParameterService{
		repository: mongorepo.NewParameterRepository(config),
	}
}

func (a *ParameterService) GetAllParameters() []ParameterViewModel {
	parameters := make([]ParameterViewModel, 0)
	entities, err := a.repository.GetAll(context.Background())
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		parameters = append(parameters, ParameterViewModel{
			Id:        entity.Id,
			AccountId: entity.AccountId,
			Region:    entity.Region,
			Key:       entity.Key,
			Value:     entity.Value,
			Version:   entity.Version,
			OperateAt: time.Unix(entity.OperateAt, 0),
		})
	}

	return parameters
}

func (a *ParameterService) GetParametersByPaging(paging core.Paging) []ParameterViewModel {
	parameters := make([]ParameterViewModel, 0)
	entities, err := a.repository.Paging(context.Background(), paging)
	if err != nil {
		return parameters
	}

	for _, entity := range entities {
		parameters = append(parameters, ParameterViewModel{
			Id:        entity.Id,
			AccountId: entity.AccountId,
			Region:    entity.Region,
			Key:       entity.Key,
			Value:     entity.Value,
			Version:   entity.Version,
			OperateAt: time.Unix(entity.OperateAt, 0),
		})
	}

	return parameters
}

func (a *ParameterService) GetParameterByID(id string) *ParameterViewModel {
	entity, err := a.repository.Get(context.Background(), id)
	if err != nil {
		return nil
	}

	parameter := &ParameterViewModel{
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
