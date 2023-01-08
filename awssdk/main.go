package main

import (
	"github.com/futugyousuzu/goproject/awsgolang/servicediscovery"
	// "github.com/futugyousuzu/goproject/awsgolang/cloudwatch"
	// "github.com/futugyousuzu/goproject/awsgolang/cloudwatchlogs"
	// "github.com/futugyousuzu/goproject/awsgolang/ecr"
)

func main() {
	// servicediscovery.ListNamespace()
	// servicediscovery.ListServices()
	// servicediscovery.RegisterInstance()
	servicediscovery.CreateService()

	// cloudwatch.GetMetricData()
	// cloudwatch.ListDashboards()
	// cloudwatch.ListMetrics()
	// cloudwatch.GetMetricStatistics()

	// cloudwatchlogs.DescribeExportTasks()
	// cloudwatchlogs.DescribeLogGroups()
	// cloudwatchlogs.GetLogEvents()
	// cloudwatchlogs.DescribeLogStreams()
	// cloudwatchlogs.GetLogGroupFields()
	// cloudwatchlogs.DescribeQueries()

	// ecr.DescribeRepositories()
	// ecr.CreateRepository()
	// ecr.DeleteRepository()
}
