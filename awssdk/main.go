package main

import (
	"fmt"

	"github.com/futugyousuzu/goproject/awsgolang/cloudwatchdemo"
	"github.com/futugyousuzu/goproject/awsgolang/servicediscoverydemo"
)

func main() {
	servicediscoverydemo.CreateNamespace()
	fmt.Println("this is separator!")
	cloudwatchdemo.GetMetricData()
}
