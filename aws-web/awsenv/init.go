package awsenv

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

var (
	EmptyContext            context.Context = context.Background()
	Cfg                     aws.Config
	CfgWithProfile          func(key string, secret string) error
	CfgWithProfileAndRegion func(key string, secret string, region string) error
	NamespaceId             string
	NamespaceName           string
	CloudMapServiceName     string
	UserName                string
	Password                string
	GroupName               string
	ECSClusterName          string
	ECSServiceName          string
	AwsConfigFilePath       string
	HostedZoneId            string
	AwsSecretName           string
)

func init() {
	NamespaceId = os.Getenv("CLOUD_MAP_NAMESPACE_ID")
	HostedZoneId = os.Getenv("ROUTE53_HOSTED_ZONEID")
	NamespaceName = os.Getenv("CLOUD_MAP_NAMESPACE")
	CloudMapServiceName = os.Getenv("CLOUD_MAP_SERVICE_NAME")
	UserName = os.Getenv("IAM_USER_NAME")
	Password = os.Getenv("IAM_USER_PASSWORD")
	GroupName = os.Getenv("IAM_GROUP_NAME")
	ECSClusterName = os.Getenv("ECS_CLUSTER_NAME")
	ECSServiceName = os.Getenv("ECS_SERVICE_NAME")
	AwsConfigFilePath = os.Getenv("AWS_CONFIG_FILE_PATH")
	AwsSecretName = os.Getenv("AWS_SECRET_NAME")

	var err error
	Cfg, err = config.LoadDefaultConfig(EmptyContext)
	if err != nil {
		log.Fatal(err)
	}

	CfgWithProfile = func(key string, secret string) error {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(key, secret, ""),
			),
		)
		if err != nil {
			log.Fatal(err)
			return err
		}

		Cfg = cfg
		return nil
	}

	CfgWithProfileAndRegion = func(key string, secret string, region string) error {
		err := CfgWithProfile(key, secret)
		if err != nil {
			return err
		}

		Cfg.Region = region
		return nil
	}
}

func DeleteAll() {
}
