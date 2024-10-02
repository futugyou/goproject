package vault_provider

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

type AWSClient struct {
	svc *ssm.Client
}

// aws ssm path like '/Finance/Prod/IAD/WinServ2016/license33'
func NewAWSClient() *AWSClient {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return &AWSClient{
		svc: ssm.NewFromConfig(cfg),
	}
}

func (s *AWSClient) Get(ctx context.Context, key string) (*ProviderVault, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	output, err := s.svc.GetParameter(ctx, input)
	if err != nil {
		return nil, err
	}
	if output.Parameter != nil {
		return &ProviderVault{
			Key:       key,
			Value:     *output.Parameter.Value,
			CreatedAt: *output.Parameter.LastModifiedDate,
		}, nil
	}
	return nil, fmt.Errorf("can not found secret with name %s in aws", key)
}

func (s *AWSClient) Search(ctx context.Context, prefix string) ([]ProviderVault, error) {
	var providerVaults = []ProviderVault{}
	var nextToken *string = nil
	for {
		input := &ssm.GetParametersByPathInput{
			Path:           aws.String(prefix),
			WithDecryption: aws.Bool(true),
			NextToken:      nextToken,
		}
		output, err := s.svc.GetParametersByPath(ctx, input)
		if err != nil {
			return nil, err
		}
		for _, v := range output.Parameters {
			if v.Value != nil {
				providerVaults = append(providerVaults, ProviderVault{
					Key:       *v.Name,
					Value:     *v.Value,
					CreatedAt: *v.LastModifiedDate,
				})
			}
		}
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return providerVaults, nil
}

func (s *AWSClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	putInput := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Overwrite: aws.Bool(true),
		Type:      types.ParameterTypeString,
	}
	_, err := s.svc.PutParameter(ctx, putInput)
	if err != nil {
		return nil, err
	}
	return &ProviderVault{
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *AWSClient) Delete(ctx context.Context, key string) error {
	input := ssm.DeleteParameterInput{
		Name: aws.String(key),
	}
	_, err := s.svc.DeleteParameter(ctx, &input)
	return err
}
