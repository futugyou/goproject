package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"github.com/futugyou/extensions"

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
	if !account.Valid {
		return nil, fmt.Errorf("account %s is expired", account.Id)
	}
	awsenv.CfgWithProfileAndRegion(account.AccessKeyId, account.SecretAccessKey, account.Region)
	svc := ssm.NewFromConfig(awsenv.Cfg)
	var result = []model.SSMData{}
	if len(filter.Name) > 0 {
		return s.getParameters(ctx, []string{filter.Name}, svc)
	} else {
		names, err := s.getParameterDescribes(ctx, svc)
		if err != nil {
			return s.getParameters(ctx, names, svc)
		}
	}

	return result, nil
}

func (*SSMService) getParameters(ctx context.Context, names []string, svc *ssm.Client) ([]model.SSMData, error) {
	var result = []model.SSMData{}

	namesList := extensions.SplitArray(names, 10)
	for _, names := range namesList {
		input := &ssm.GetParametersInput{
			Names:          names,
			WithDecryption: aws.Bool(true),
		}

		parameters, err := svc.GetParameters(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, p := range parameters.Parameters {
			result = append(result, model.SSMData{
				Key:        *p.Name,
				Value:      *p.Value,
				CreateDate: *p.LastModifiedDate,
			})
		}
	}

	return result, nil
}

func (*SSMService) getParameterDescribes(ctx context.Context, svc *ssm.Client) ([]string, error) {
	var result = []string{}
	var nextToken *string = nil

	for {
		input := &ssm.DescribeParametersInput{
			NextToken:  nextToken,
			MaxResults: aws.Int32(50),
		}
		output, err := svc.DescribeParameters(ctx, input)
		if err != nil {
			return nil, err
		}
		nextToken = output.NextToken
		for _, p := range output.Parameters {
			if p.Name != nil {
				result = append(result, *p.Name)
			}
		}
		if nextToken == nil {
			break
		}
	}
	return result, nil
}
