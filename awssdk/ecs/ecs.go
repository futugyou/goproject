package ecs

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *ecs.Client
)

func init() {
	svc = ecs.NewFromConfig(awsenv.Cfg)
}

func DescribeClusters() {
	input := &ecs.ListClustersInput{}
	output, err := svc.ListClusters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	describeInput := &ecs.DescribeClustersInput{
		Clusters: output.ClusterArns,
	}
	describeOutput, err := svc.DescribeClusters(awsenv.EmptyContext, describeInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, cluster := range describeOutput.Clusters {
		fmt.Println("ClusterName:", *cluster.ClusterName, "\tArn:", *cluster.ClusterArn, "\tStatus:", *cluster.Status)
	}
}
