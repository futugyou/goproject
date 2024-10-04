package vault_provider

import (
	"context"
	"fmt"
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
func NewAWSClient() (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}
	return &AWSClient{
		svc: ssm.NewFromConfig(cfg),
	}, nil
}

func (s *AWSClient) Search(ctx context.Context, key string) (*ProviderVault, error) {
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

func (s *AWSClient) PrefixSearch(ctx context.Context, prefix string) (map[string]ProviderVault, error) {
	var providerVaults = map[string]ProviderVault{}
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
				providerVaults[*v.Name] = ProviderVault{
					Key:       *v.Name,
					Value:     *v.Value,
					CreatedAt: *v.LastModifiedDate,
				}
			}
		}
		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}
	return providerVaults, nil
}

func (s *AWSClient) BatchSearch(ctx context.Context, keys []string) (map[string]ProviderVault, error) {
	var providerVaults = map[string]ProviderVault{}
	input := &ssm.GetParametersInput{
		Names:          keys,
		WithDecryption: aws.Bool(true),
	}
	output, err := s.svc.GetParameters(ctx, input)
	if err != nil {
		return nil, err
	}
	for _, parameter := range output.Parameters {
		providerVaults[*parameter.Name] = ProviderVault{
			Key:       *parameter.Name,
			Value:     *parameter.Value,
			CreatedAt: *parameter.LastModifiedDate,
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

	if _, err := s.svc.PutParameter(ctx, putInput); err != nil {
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

func (s *AWSClient) SearchAsync(ctx context.Context, key string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Search(ctx, key)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AWSClient) PrefixSearchAsync(ctx context.Context, prefix string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.PrefixSearch(ctx, prefix)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AWSClient) BatchSearchAsync(ctx context.Context, keys []string) (<-chan map[string]ProviderVault, <-chan error) {
	resultChan := make(chan map[string]ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.BatchSearch(ctx, keys)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AWSClient) UpsertAsync(ctx context.Context, key string, value string) (<-chan *ProviderVault, <-chan error) {
	resultChan := make(chan *ProviderVault, 1)
	errorChan := make(chan error, 1)
	go func() {
		defer close(resultChan)
		defer close(errorChan)
		result, err := s.Upsert(ctx, key, value)
		if err != nil {
			errorChan <- err
			return
		}
		resultChan <- result
	}()
	return resultChan, errorChan
}

func (s *AWSClient) DeleteAsync(ctx context.Context, key string) <-chan error {
	errorChan := make(chan error, 1)
	go func() {
		defer close(errorChan)
		errorChan <- s.Delete(ctx, key)
	}()
	return errorChan
}
