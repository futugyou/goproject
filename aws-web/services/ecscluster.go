package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
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
			Id:             entity.Id,
			ClusterName:    entity.Cluster,
			ClusterArn:     entity.ClusterArn,
			ServiceName:    entity.ServiceName,
			ServiceNameArn: entity.ServiceNameArn,
			RoleArn:        entity.RoleArn,
			AccountAlias:   alias,
			OperateAt:      entity.OperateAt,
		}

		result = append(result, e)
	}
	return result, nil
}

func (e *EcsClusterService) GetServiceDetailById(id string) (*model.EcsClusterDetailViewModel, error) {
	// 1 data from mongo db
	entity, err := e.repository.GetByObjectId(context.Background(), id)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	accountService := NewAccountService()
	account := accountService.GetAccountByID(entity.AccountId)

	result := &model.EcsClusterDetailViewModel{}
	result.AccountAlias = account.Alias
	result.Id = entity.Id
	result.ClusterArn = entity.ClusterArn
	result.ClusterName = entity.Cluster
	result.OperateAt = entity.OperateAt
	result.Service = entity.ServiceName
	result.ServiceArn = entity.ServiceNameArn
	result.RoleArn = entity.RoleArn

	// 2 data from aws cloud
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	svc := ecs.NewFromConfig(awsenv.Cfg)

	describeInput := &ecs.DescribeServicesInput{
		Cluster:  aws.String(entity.Cluster),
		Services: []string{entity.ServiceName},
	}

	describeOutput, err := svc.DescribeServices(awsenv.EmptyContext, describeInput)
	if err != nil || len(describeOutput.Services) == 0 {
		return result, nil
	}

	service := describeOutput.Services[0]
	if service.NetworkConfiguration != nil && service.NetworkConfiguration.AwsvpcConfiguration != nil {
		result.SecurityGroups = service.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups
		result.Subnets = service.NetworkConfiguration.AwsvpcConfiguration.Subnets
	}

	loadBalancers := make([]string, 0)
	fmt.Println(service.LoadBalancers)
	for _, lb := range service.LoadBalancers {
		if lb.TargetGroupArn != nil {
			loadBalancers = append(loadBalancers, *lb.TargetGroupArn)
		}
	}
	result.LoadBalancers = loadBalancers

	serviceRegistries := make([]string, 0)
	for _, sr := range service.ServiceRegistries {
		if sr.RegistryArn != nil {
			serviceRegistries = append(serviceRegistries, *sr.RegistryArn)
		}
	}
	result.ServiceRegistries = serviceRegistries

	listTaskInput := &ecs.ListTaskDefinitionsInput{
		MaxResults:   aws.Int32(10),
		FamilyPrefix: aws.String(entity.ServiceName),
		Sort:         types.SortOrderDesc,
	}

	listTaskOutput, err := svc.ListTaskDefinitions(awsenv.EmptyContext, listTaskInput)
	if err != nil {
		log.Println(err)
		return result, nil
	}

	result.TaskDefinitions = listTaskOutput.TaskDefinitionArns

	return result, nil
}

func (e *EcsClusterService) CompareTaskDefinitions(compare model.EcsTaskCompare) ([]string, error) {
	result := make([]string, 0)
	entity, err := e.repository.GetByObjectId(context.Background(), compare.Id)
	if err != nil {
		log.Println(err)
		return result, err
	}

	accountService := NewAccountService()
	account := accountService.GetAccountByID(entity.AccountId)
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	data1 := describeTaskDefinition(compare.SourceTaskArn)
	data2 := describeTaskDefinition(compare.DestTaskArn)
	result = append(result, data1, data2)
	return result, nil
}

func describeTaskDefinition(taskArn string) string {
	svc := ecs.NewFromConfig(awsenv.Cfg)
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(taskArn),
	}

	output, err := svc.DescribeTaskDefinition(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err)
		return ""
	}

	data, _ := json.Marshal(output.TaskDefinition)
	return string(data)
}
