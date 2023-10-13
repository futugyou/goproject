package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
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
