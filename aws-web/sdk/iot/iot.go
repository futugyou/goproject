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
	input := &iot.ListJobsInput{}

	result, err := svc.ListJobs(awsenv.EmptyContext, input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	for _, job := range result.Jobs {
		log.Println(*job.JobId, job.Status, *job.ThingGroupId, job.TargetSelection)
	}
}
