package loadbalancing

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
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

func DescribeDescribeListeners() {
	input := &elasticloadbalancingv2.DescribeListenersInput{}

	result, err := svc.DescribeListeners(context.TODO(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, target := range result.Listeners {
		fmt.Println(*target.ListenerArn, *target.LoadBalancerArn)
		for _, act := range target.DefaultActions {
			fmt.Println(act.TargetGroupArn)
		}
	}
}

func GetTargetGroups() []types.TargetGroup {
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{}

	result, err := svc.DescribeTargetGroups(context.TODO(), input)
	if err != nil {
		log.Println(err)
		return []types.TargetGroup{}
	}

	return result.TargetGroups
}

func GetLoadbalanceListeners() []types.Listener {
	input := &elasticloadbalancingv2.DescribeListenersInput{}

	result, err := svc.DescribeListeners(context.TODO(), input)
	if err != nil {
		log.Println(err)
		return []types.Listener{}
	}

	return result.Listeners
}
