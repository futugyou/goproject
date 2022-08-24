package cloudwatchdemo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func GetMetricData() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	d := 6 * 60 * time.Minute
	svc := cloudwatch.New(sess)
	queryid := "q1"
	epression := "SELECT AVG(\"http.server.duration\") FROM \"ECS/APIGateway\" GROUP BY \"http.target\", \"service.name\""
	query := make([]*cloudwatch.MetricDataQuery, 0)
	query = append(query, &cloudwatch.MetricDataQuery{
		Id:         aws.String(queryid),
		Expression: aws.String(epression),
		Period:     aws.Int64(300),
	})

	input := &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(time.Now().UTC().Add(-d)),
		EndTime:           aws.Time(time.Now().UTC()),
		MetricDataQueries: query,
	}
	result, err := svc.GetMetricData(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case cloudwatch.ErrCodeInternalServiceFault:
				fmt.Println(cloudwatch.ErrCodeInternalServiceFault, aerr.Error())
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
