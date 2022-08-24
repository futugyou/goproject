package cloudwatchdemo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

func DescribeExportTasks() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)

	input := &cloudwatchlogs.DescribeExportTasksInput{
		Limit: aws.Int64(50),
	}
	result, err := svc.DescribeExportTasks(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}

func DescribeLogGroups() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)

	input := &cloudwatchlogs.DescribeLogGroupsInput{
		Limit: aws.Int64(50),
	}
	result, err := svc.DescribeLogGroups(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}

func GetLogEvents() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)
	d := 6 * 60 * time.Minute
	input := &cloudwatchlogs.GetLogEventsInput{
		Limit:         aws.Int64(50),
		StartTime:     aws.Int64(time.Now().UTC().Add(-d).Unix()),
		EndTime:       aws.Int64(time.Now().UTC().Unix()),
		LogGroupName:  aws.String("/ecs/oneapp-apigateway"),
		LogStreamName: aws.String("ecs/oneapp-apigateway/7c04a4d0ccd74c9cbfc7a5e6e7d76501"),
	}
	result, err := svc.GetLogEvents(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}

func DescribeLogStreams() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)
	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:        aws.Int64(50),
		LogGroupName: aws.String("/ecs/oneapp-apigateway"),
	}
	result, err := svc.DescribeLogStreams(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}

func GetLogGroupFields() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)
	d := 60 * time.Minute
	input := &cloudwatchlogs.GetLogGroupFieldsInput{
		LogGroupName: aws.String("/ecs/oneapp-apigateway"),
		Time:         aws.Int64(time.Now().UTC().Add(-d).Unix()),
	}
	result, err := svc.GetLogGroupFields(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}

func DescribeQueries() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := cloudwatchlogs.New(sess)
	input := &cloudwatchlogs.DescribeQueriesInput{
		LogGroupName: aws.String("/ecs/oneapp-apigateway"),
	}
	result, err := svc.DescribeQueries(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatchlogs.ErrCodeServiceUnavailableException:
				fmt.Println(cloudwatchlogs.ErrCodeServiceUnavailableException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
	}
	fmt.Println(result)
}
