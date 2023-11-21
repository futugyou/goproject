package iot

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *iot.Client
)

func init() {
	svc = iot.NewFromConfig(awsenv.Cfg)
}

func ListJobs() {
	var nextToken *string = nil
	for {
		input := &iot.ListJobsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListJobs(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, job := range result.Jobs {
			log.Println(*job.JobId, job.Status, *job.ThingGroupId, job.TargetSelection)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThings() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingsInput{
			NextToken: nextToken,
		}

		result, err := svc.ListThings(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, thing := range result.Things {
			log.Println(*thing.ThingArn, *thing.ThingName, *thing.ThingTypeName, thing.Version)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingTypes() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingTypesInput{
			NextToken: nextToken,
		}

		result, err := svc.ListThingTypes(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, ty := range result.ThingTypes {
			log.Println(*ty.ThingTypeName, *ty.ThingTypeArn)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingGroups() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingGroupsInput{
			NextToken: nextToken,
			Recursive:aws.Bool(true),
		}

		result, err := svc.ListThingGroups(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, group := range result.ThingGroups {
			log.Println(*group.GroupArn,*group.GroupName)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListThingRegistrationTasks() {
	var nextToken *string = nil
	for {
		input := &iot.ListThingRegistrationTasksInput{
			NextToken: nextToken, 
		}

		result, err := svc.ListThingRegistrationTasks(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, taskid := range result.TaskIds {
			log.Println(taskid)
		}
		if result.NextToken == nil {
			return
		}
	}
}


func ListTopicRuleDestinations() {
	var nextToken *string = nil
	for {
		input := &iot.ListTopicRuleDestinationsInput{
			NextToken: nextToken, 
		}

		result, err := svc.ListTopicRuleDestinations(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, summ := range result.DestinationSummaries {
			log.Println(*summ.Arn,summ.Status,*summ.StatusReason,summ.HttpUrlSummary,summ.VpcDestinationSummary)
		}
		if result.NextToken == nil {
			return
		}
	}
}

func ListTopicRules() {
	var nextToken *string = nil
	for {
		input := &iot.ListTopicRulesInput{
			NextToken: nextToken, 
		}

		result, err := svc.ListTopicRules(awsenv.EmptyContext, input)
		if err != nil {
			log.Println(err.Error())
			return
		}
		nextToken = result.NextToken

		for _, rule := range result.Rules {
			log.Println(*rule.RuleArn,*rule.RuleDisabled,*rule.RuleName,*rule.TopicPattern,*rule.CreatedAt)
		}
		if result.NextToken == nil {
			return
		}
	}
}
