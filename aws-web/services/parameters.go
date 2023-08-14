package services

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/entity"
	"github.com/futugyousuzu/goproject/awsgolang/repository"
	"github.com/futugyousuzu/goproject/awsgolang/repository/mongorepo"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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

func (a *ParameterService) SyncAllParameter() {
	log.Println("start..")
	accountService := NewAccountService()
	accounts := accountService.GetAllAccounts()

	entities := make([]entity.ParameterEntity, 0)
	logs := make([]entity.ParameterLogEntity, 0)
	for _, account := range accounts {
		awsenv.CfgForVercelWithRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
		parameters, err := a.getAllParametersFromAWS()
		if err != nil {
			continue
		}

		names := make([]string, len(parameters))
		for i := 0; i < len(parameters); i++ {
			names[i] = *parameters[i].Name
		}

		details, err := a.getParametersDatail(names)
		if err != nil {
			continue
		}

		for _, d := range details {
			p := entity.ParameterEntity{
				AccountId: account.Id,
				Region:    account.Region,
				Key:       *d.Name,
				Value:     *d.Value,
				Version:   strconv.FormatInt(d.Version, 10),
				OperateAt: time.Now().Unix(),
			}

			entities = append(entities, p)

			l := entity.ParameterLogEntity{
				AccountId: account.Id,
				Region:    account.Region,
				Key:       *d.Name,
				Value:     *d.Value,
				Version:   strconv.FormatInt(d.Version, 10),
				OperateAt: time.Now().Unix(),
			}

			logs = append(logs, l)
		}
	}

	log.Println("get finish, count: ", len(entities))
	err := a.repository.BulkWrite(context.Background(), entities)
	log.Println("parameter write finish: ", err)
	err = a.logRepository.BulkWrite(context.Background(), logs)
	log.Println("log write finish: ", err)
}

func (a *ParameterService) getAllParametersFromAWS() ([]types.ParameterMetadata, error) {
	svc := ssm.NewFromConfig(awsenv.Cfg)
	totals := make([]types.ParameterMetadata, 0)

	var nextToken *string = nil
	for {
		var input *ssm.DescribeParametersInput
		if nextToken == nil {
			input = &ssm.DescribeParametersInput{
				MaxResults: aws.Int32(50), // max value 50
			}
		} else {
			input = &ssm.DescribeParametersInput{
				MaxResults: aws.Int32(50), // max value 50
				NextToken:  nextToken,
			}
		}

		output, err := svc.DescribeParameters(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err)
			break
		}

		nextToken = output.NextToken
		if len(output.Parameters) == 0 {
			log.Println(err)
			break
		}

		totals = append(totals, output.Parameters...)

		if nextToken == nil {
			break
		}
	}

	return totals, nil
}

func (a *ParameterService) getParametersDatail(names []string) ([]types.Parameter, error) {
	svc := ssm.NewFromConfig(awsenv.Cfg)

	totals := make([]types.Parameter, 0)

	for {
		if len(names) == 0 {
			break
		}

		t := names
		if len(t) > 10 {
			t = names[:10]
		}

		input := &ssm.GetParametersInput{
			Names:          t,
			WithDecryption: aws.Bool(true),
		}

		output, err := svc.GetParameters(awsenv.EmptyContext, input)
		if len(names) > 10 {
			names = names[10:]
		} else {
			names = []string{}
		}

		if err != nil {
			log.Println(err)
			continue
		}

		if len(output.Parameters) == 0 {
			continue
		}

		totals = append(totals, output.Parameters...)

	}

	return totals, nil
}
