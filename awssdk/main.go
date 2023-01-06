package main

import (
	"github.com/futugyousuzu/goproject/awsgolang/servicediscovery"

	"github.com/futugyousuzu/goproject/awsgolang/cloudwatchlogs"
)

func main() {
	servicediscovery.ListNamespace()
	// cloudwatchdemo.GetMetricData()
	// cloudwatchdemo.GetDashboard()
	// cloudwatchdemo.ListMetrics()
	// cloudwatchdemo.GetMetricStatistics()
	cloudwatchlogs.DescribeExportTasks()
	cloudwatchlogs.DescribeLogGroups()
	cloudwatchlogs.GetLogEvents()
	cloudwatchlogs.DescribeLogStreams()
	cloudwatchlogs.GetLogGroupFields()
	cloudwatchlogs.DescribeQueries()
}
