package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

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
			Value:     entity.Value,
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
			Value:     entity.Value,
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

func (a *ParameterService) getParametersDatail(names []string) ([]types.Parameter, error) {
	svc := ssm.NewFromConfig(awsenv.Cfg)
	input := &ssm.GetParametersInput{
		Names:          names,
		WithDecryption: aws.Bool(true),
	}

	output, err := svc.GetParameters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if len(output.Parameters) == 0 {
		fmt.Println("no data found")
		return nil, fmt.Errorf("no data found")
	}

	return output.Parameters, nil
}
