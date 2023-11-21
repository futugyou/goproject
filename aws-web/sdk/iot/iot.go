package iot

import (
	"log"

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
