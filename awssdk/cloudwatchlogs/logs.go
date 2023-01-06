package cloudwatchlogs

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *cloudwatchlogs.Client
)

func init() {
	svc = cloudwatchlogs.NewFromConfig(awsenv.Cfg)
}

func DescribeExportTasks() {
	input := &cloudwatchlogs.DescribeExportTasksInput{
		Limit: aws.Int32(50),
	}
	result, err := svc.DescribeExportTasks(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, task := range result.ExportTasks {
		fmt.Println("TaskId:", *task.TaskId, "\tTaskName:", *task.TaskName, "\tLogGroupName:", *task.LogGroupName)
	}
}

func DescribeLogGroups() {
	// d := 6 * 60 * time.Minute
	input := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int32(10),
	}
	result, err := svc.DescribeLogGroups(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, group := range result.LogGroups {
		fmt.Println("LogGroupName:", *group.LogGroupName)
		input := &cloudwatchlogs.DescribeLogStreamsInput{
			Limit:        aws.Int32(10),
			LogGroupName: group.LogGroupName,
		}
		result, err := svc.DescribeLogStreams(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		for _, stream := range result.LogStreams {
			fmt.Println("\tLogStreamName:", *stream.LogStreamName)
			input := &cloudwatchlogs.GetLogEventsInput{
				Limit: aws.Int32(10),
				// StartTime:     aws.Int64(time.Now().UTC().Add(-d).Unix()),
				// EndTime:       aws.Int64(time.Now().UTC().Unix()),
				LogGroupName:  group.LogGroupName,
				LogStreamName: stream.LogStreamName,
			}
			result, err := svc.GetLogEvents(awsenv.EmptyContext, input)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			for _, event := range result.Events {
				fmt.Println("\t\tMessage:", *event.Message)
			}
		}
	}
}

// func DescribeLogStreams() {
// 	input := &cloudwatchlogs.DescribeLogStreamsInput{
// 		Limit:        aws.Int32(50),
// 		LogGroupName: aws.String("/eks/openTelemetry"),
// 	}
// 	result, err := svc.DescribeLogStreams(awsenv.EmptyContext, input)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(result)
// }

// func GetLogEvents() {
// 	d := 6 * 60 * time.Minute
// 	input := &cloudwatchlogs.GetLogEventsInput{
// 		Limit:         aws.Int32(50),
// 		StartTime:     aws.Int64(time.Now().UTC().Add(-d).Unix()),
// 		EndTime:       aws.Int64(time.Now().UTC().Unix()),
// 		LogGroupName:  aws.String("/eks/openTelemetry"),
// 		LogStreamName: aws.String("/eks/openTelemetry/7c04a4d0ccd74c9cbfc7a5e6e7d76501"),
// 	}
// 	result, err := svc.GetLogEvents(awsenv.EmptyContext, input)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println(result)
// }

func GetLogGroupFields() {
	d := 60 * time.Minute
	input := &cloudwatchlogs.GetLogGroupFieldsInput{
		LogGroupName: aws.String("/eks/openTelemetry"),
		Time:         aws.Int64(time.Now().UTC().Add(-d).Unix()),
	}
	result, err := svc.GetLogGroupFields(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}

func DescribeQueries() {
	input := &cloudwatchlogs.DescribeQueriesInput{
		LogGroupName: aws.String("/eks/openTelemetry"),
	}
	result, err := svc.DescribeQueries(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
