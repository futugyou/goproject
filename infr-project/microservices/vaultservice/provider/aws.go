package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/futugyou/vaultservice/options"
)

type AWSClient struct {
	svc *ssm.Client
}

// aws ssm path like '/Finance/Prod/IAD/WinServ2016/license33'
func NewAWSClient(opts *options.Options) (*AWSClient, error) {
	var cfg aws.Config
	var err error
	if len(opts.AwsAccessKeyID) > 0 && len(opts.AwsSecretAccessKey) > 0 {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(opts.AwsAccessKeyID, opts.AwsSecretAccessKey, ""),
			),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.Background())
	}

	if err != nil {
		return nil, err
	}

	if len(opts.AwsRegion) > 0 {
		cfg.Region = opts.AwsRegion
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
