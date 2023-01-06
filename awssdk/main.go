package main

import (
	"github.com/futugyousuzu/goproject/awsgolang/servicediscovery"

	"github.com/futugyousuzu/goproject/awsgolang/cloudwatch"
	"github.com/futugyousuzu/goproject/awsgolang/cloudwatchlogs"
)

func main() {
	servicediscovery.ListNamespace()
	cloudwatch.GetMetricData()
	cloudwatch.GetDashboard()
	cloudwatch.ListMetrics()
	cloudwatch.GetMetricStatistics()
	cloudwatchlogs.DescribeExportTasks()
	cloudwatchlogs.DescribeLogGroups()
	cloudwatchlogs.GetLogEvents()
	cloudwatchlogs.DescribeLogStreams()
	cloudwatchlogs.GetLogGroupFields()
	cloudwatchlogs.DescribeQueries()
}
