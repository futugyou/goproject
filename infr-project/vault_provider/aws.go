package vault_provider

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
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
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(prefix),
		WithDecryption: aws.Bool(true),
	}
	//TODO: handle NextToken
	output, err := s.svc.GetParametersByPath(ctx, input)
	if err != nil {
		return nil, err
	}
	var providerVaults = []ProviderVault{}
	for _, v := range output.Parameters {
		if v.Value != nil {
			providerVaults = append(providerVaults, ProviderVault{
				Key:       *v.Name,
				Value:     *v.Value,
				CreatedAt: *v.LastModifiedDate,
			})
		}
	}
	return providerVaults, nil
}

func (s *AWSClient) Upsert(ctx context.Context, key string, value string) (*ProviderVault, error) {
	return nil, nil
}

func (s *AWSClient) Delete(ctx context.Context, key string) error {
	return nil
}
