package cloudwatch

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *cloudwatch.Client
)

func init() {
	svc = cloudwatch.NewFromConfig(awsenv.Cfg)
}

func GetMetricData() {
	d := 6 * 60 * time.Minute
	queryid := "q1"
	epression := "SELECT AVG(\"http.server.duration\") FROM \"ECS/APIGateway\" GROUP BY \"http.target\", \"service.name\""
	query := make([]types.MetricDataQuery, 0)
	query = append(query, types.MetricDataQuery{
		Id:         aws.String(queryid),
		Expression: aws.String(epression),
		Period:     aws.Int32(300),
	})

	input := &cloudwatch.GetMetricDataInput{
		StartTime:         aws.Time(time.Now().UTC().Add(-d)),
		EndTime:           aws.Time(time.Now().UTC()),
		MetricDataQueries: query,
	}
	result, err := svc.GetMetricData(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, data := range result.MetricDataResults {
		fmt.Println("label:", *data.Label)
	}

	for _, message := range result.Messages {
		fmt.Println("code:", *message.Code, "\tmessage:", *message.Value)
	}
}

func ListDashboards() {
	input := &cloudwatch.ListDashboardsInput{}
	result, err := svc.ListDashboards(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, entry := range result.DashboardEntries {
		// fmt.Println("DashboardName:", *entry.DashboardName)
		input := &cloudwatch.GetDashboardInput{
			DashboardName: entry.DashboardName,
		}
		result, err := svc.GetDashboard(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fmt.Println("DashboardName:", *result.DashboardName, "\tDashboardBody:", *result.DashboardBody)
	}
}

func ListMetrics() {
	input := &cloudwatch.ListMetricsInput{}
	result, err := svc.ListMetrics(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, metric := range result.Metrics {
		fmt.Print("Namespace:", *metric.Namespace, "\tMetricName:", *metric.MetricName, "\tDimensions:")
		for _, dimension := range metric.Dimensions {
			fmt.Print(*dimension.Name, *dimension.Value)
		}
		fmt.Println()
	}
}

func GetMetricStatistics() {
	d := 6 * 60 * time.Minute
	statistics := make([]types.Statistic, 0)
	statistics = append(statistics, types.StatisticSum)
	input := &cloudwatch.GetMetricStatisticsInput{
		StartTime:  aws.Time(time.Now().UTC().Add(-d)),
		EndTime:    aws.Time(time.Now().UTC()),
		Period:     aws.Int32(300),
		MetricName: aws.String("ClusterName,ServiceName"),
		Namespace:  aws.String("ECS"),
		Statistics: statistics,
	}
	result, err := svc.GetMetricStatistics(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("label:", *result.Label)
	for _, point := range result.Datapoints {
		fmt.Print("Average:", point.Average, "\tSampleCount:", point.SampleCount, "\tExtendedStatistics:")
		for k, v := range point.ExtendedStatistics {
			fmt.Print(k, v)
		}
		fmt.Println()
	}
}
