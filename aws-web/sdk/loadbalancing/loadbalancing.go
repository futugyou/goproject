package loadbalancing

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	svc = elasticloadbalancingv2.NewFromConfig(awsenv.Cfg)
	input := &elasticloadbalancingv2.DescribeTargetGroupsInput{}

	result, err := svc.DescribeTargetGroups(context.TODO(), input)
	if err != nil {
		log.Println("describe target group error")
		return []types.TargetGroup{}
	}

	return result.TargetGroups
}

func GetLoadbalanceListeners(lbs []string) []types.Listener {
	svc = elasticloadbalancingv2.NewFromConfig(awsenv.Cfg)
	listeners := make([]types.Listener, 0)
	for _, ls := range lbs {
		input := &elasticloadbalancingv2.DescribeListenersInput{
			LoadBalancerArn: aws.String(ls),
		}

		result, err := svc.DescribeListeners(context.TODO(), input)
		if err != nil {
			log.Println("describe loadbalance listener error")
			continue
		}
		listeners = append(listeners, result.Listeners...)
	}

	return listeners
}
