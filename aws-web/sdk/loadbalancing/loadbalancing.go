package loadbalancing

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *elasticloadbalancingv2.Client
)

func init() {
	svc = elasticloadbalancingv2.NewFromConfig(awsenv.Cfg)
}

func DescribeTargetGroups() {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{}

	result, err := svc.DescribeTargetGroups(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, target := range result.TargetGroups {
		fmt.Println(*target.TargetGroupArn, *target.TargetGroupName, *target.VpcId)
	}
}
