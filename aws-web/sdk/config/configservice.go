package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *configservice.Client
)

func init() {
	svc = configservice.NewFromConfig(awsenv.Cfg)
}

func DeliverConfigSnapshot() {
	input := &configservice.DeliverConfigSnapshotInput{
		DeliveryChannelName: aws.String("default"),
	}
	result, err := svc.DeliverConfigSnapshot(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("ConfigSnapshotId:", *result.ConfigSnapshotId)
}

func DescribeConfigRules() {
	var nextToken *string = nil

	for {
		input := &configservice.DescribeConfigRulesInput{
			NextToken: nextToken,
		}

		result, err := svc.DescribeConfigRules(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		nextToken = result.NextToken

		for _, v := range result.ConfigRules {
			fmt.Println(*v.ConfigRuleName)
			inputC := &configservice.DeleteConfigRuleInput{
				ConfigRuleName: v.ConfigRuleName,
			}
			svc.DeleteConfigRule(awsenv.EmptyContext, inputC)
		}

		if result.NextToken == nil {
			return
		}
	}
}

func DescribeConfigurationRecorders() {
	input := &configservice.DescribeConfigurationRecordersInput{}
	result, err := svc.DescribeConfigurationRecorders(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range result.ConfigurationRecorders {
		fmt.Println(*v.Name, *v.RoleARN)
	}
}

func DeleteConfigurationRecorder() {
	input := &configservice.DeleteConfigurationRecorderInput{
		ConfigurationRecorderName: aws.String("default"),
	}
	_, err := svc.DeleteConfigurationRecorder(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func StartConfigurationRecorder() {
	input := &configservice.StartConfigurationRecorderInput{
		ConfigurationRecorderName: aws.String("default"),
	}

	_, err := svc.StartConfigurationRecorder(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("StartConfigurationRecorder ok")
}

func PutConfigurationRecorder() {
	input := &configservice.PutConfigurationRecorderInput{
		ConfigurationRecorder: &types.ConfigurationRecorder{
			Name:    aws.String("default"),
			RoleARN: aws.String(" "), //roles: AWSServiceRoleForConfig, Policy: AWSConfigServiceRolePolicy
		},
	}

	_, err := svc.PutConfigurationRecorder(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("PutConfigurationRecorder ok")
}

func PutDeliveryChannel() {
	input := &configservice.PutDeliveryChannelInput{
		DeliveryChannel: &types.DeliveryChannel{
			Name:         aws.String("default"),
			S3BucketName: aws.String("s3_name"),
		},
	}

	_, err := svc.PutDeliveryChannel(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("PutDeliveryChannel ok")
}

func GetAwsConfigData() {
	ctx := context.Background()
	var nextToken *string = nil
	results := make([]string, 0)
	for {
		input := &configservice.SelectResourceConfigInput{
			Expression: aws.String(`
	SELECT
		version,
		accountId,
		configurationItemCaptureTime,
		configurationItemStatus,
		configurationStateId,
		arn,
		resourceType,
		resourceId,
		resourceName,
		awsRegion,
		availabilityZone,
		tags,
		relatedEvents,
		relationships,
		configuration,
		supplementaryConfiguration,
		resourceTransitionStatus,
		resourceCreationTime
	WHERE
		resourceType <> 'AWS::Backup::RecoveryPoint'
		and resourceType <> 'AWS::CodeDeploy::DeploymentConfig'
		and resourceType <> 'AWS::RDS::DBSnapshot'
  `),
			Limit:     100,
			NextToken: nextToken,
		}
		fmt.Println(1)
		output, err := svc.SelectResourceConfig(ctx, input)
		if err != nil {
			fmt.Println("select aws config resource error")
			return
		}
		fmt.Println(2)
		results = append(results, output.Results...)
		nextToken = output.NextToken

		if output.NextToken == nil {
			break
		}
	}
	fmt.Println(results)
}
