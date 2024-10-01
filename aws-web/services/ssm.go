package services

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	model "github.com/futugyousuzu/goproject/awsgolang/viewmodel"
)

type SSMService struct {
}

func NewSSMService() *SSMService {
	return &SSMService{}
}

func (s *SSMService) SearchSSMData(ctx context.Context, filter model.SSMDataFilter) ([]model.SSMData, error) {
	accountService := NewAccountService()
	account := accountService.GetAccountByID(ctx, filter.AccountId)
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	svc := ssm.NewFromConfig(awsenv.Cfg)

	input := &ssm.GetParametersInput{
		Names:          []string{filter.Name}, // max count 10, but we have only one
		WithDecryption: aws.Bool(true),
	}

	parameters, err := svc.GetParameters(awsenv.EmptyContext, input)
	if err != nil {
		return nil, err
	}

	var result = []model.SSMData{}
	for _, p := range parameters.Parameters {
		result = append(result, model.SSMData{
			Key:        *p.Name,
			Value:      *p.Value,
			CreateDate: *p.LastModifiedDate,
		})
	}

	return result, nil
}
