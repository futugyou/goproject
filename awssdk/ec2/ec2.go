package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
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
	}
}
