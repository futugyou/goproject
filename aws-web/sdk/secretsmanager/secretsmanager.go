package secretsmanager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *secretsmanager.Client
)

func init() {
	svc = secretsmanager.NewFromConfig(awsenv.Cfg)
}

func GetSecretValue() {
	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(awsenv.AwsSecretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decrypts secret using the associated KMS key.
	var secretString string = *result.SecretString
	fmt.Println(secretString)
}
