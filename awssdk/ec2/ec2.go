package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *ec2.Client
)

func init() {
	svc = ec2.NewFromConfig(awsenv.Cfg)
}

func DescribeSecurityGroups() {
	input := ec2.DescribeSecurityGroupsInput{}
	output, err := svc.DescribeSecurityGroups(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, sg := range output.SecurityGroups {
		fmt.Print(*sg.GroupName, *sg.GroupId)
		for _, permission := range sg.IpPermissions {
			fmt.Println("\t", *permission.IpProtocol, permission.FromPort, permission.ToPort)
			for _, pair := range permission.UserIdGroupPairs {
				fmt.Println("\t", pair.GroupName)
			}
		}
	}
}

func DescribeVpcs() {
	input := ec2.DescribeVpcsInput{}
	output, err := svc.DescribeVpcs(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, vpc := range output.Vpcs {
		fmt.Println(*vpc.VpcId, *vpc.CidrBlock, *vpc.DhcpOptionsId, vpc.InstanceTenancy, *vpc.IsDefault, *vpc.OwnerId, vpc.State)
		for _, set := range vpc.CidrBlockAssociationSet {
			fmt.Println("\t", *set.AssociationId, *set.CidrBlock, set.CidrBlockState.State)
		}

		// 1 VpcAttributeNameEnableDnsSupport
		input := ec2.DescribeVpcAttributeInput{
			VpcId:     vpc.VpcId,
			Attribute: types.VpcAttributeNameEnableDnsSupport,
		}
		output, err := svc.DescribeVpcAttribute(awsenv.EmptyContext, &input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tEnableDnsSupport:", *output.EnableDnsSupport.Value)

		// 2 VpcAttributeNameEnableDnsHostnames
		input = ec2.DescribeVpcAttributeInput{
			VpcId:     vpc.VpcId,
			Attribute: types.VpcAttributeNameEnableDnsHostnames,
		}
		output, err = svc.DescribeVpcAttribute(awsenv.EmptyContext, &input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tEnableDnsHostnames:", *output.EnableDnsHostnames.Value)

		// 3 VpcAttributeNameEnableNetworkAddressUsageMetrics
		input = ec2.DescribeVpcAttributeInput{
			VpcId:     vpc.VpcId,
			Attribute: types.VpcAttributeNameEnableNetworkAddressUsageMetrics,
		}
		output, err = svc.DescribeVpcAttribute(awsenv.EmptyContext, &input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tEnableNetworkAddressUsageMetrics:", *output.EnableNetworkAddressUsageMetrics.Value)

		fmt.Println()
	}
}
