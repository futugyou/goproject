package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
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

func CreateVpc() {
	input := ec2.CreateVpcInput{
		CidrBlock:       aws.String("10.0.0.0/16"), // this must >=16 and <=28
		InstanceTenancy: types.TenancyDedicated,
		TagSpecifications: []types.TagSpecification{
			{
				ResourceType: types.ResourceTypeVpc,
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("Terraform-Vpc-01"),
					},
				},
			},
		},
	}
	output, err := svc.CreateVpc(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	vpc := output.Vpc
	fmt.Println(*vpc.VpcId, *vpc.CidrBlock, *vpc.DhcpOptionsId, vpc.InstanceTenancy, *vpc.IsDefault, *vpc.OwnerId, vpc.State)
	for _, set := range vpc.CidrBlockAssociationSet {
		fmt.Println("\t", *set.AssociationId, *set.CidrBlock, set.CidrBlockState.State)
	}
}

func AssociateVpcCidrBlock() {
	input := ec2.AssociateVpcCidrBlockInput{
		VpcId:     aws.String("vpc-0e6c30199bc54494b"),
		CidrBlock: aws.String("10.1.0.0/16"),
	}
	output, err := svc.AssociateVpcCidrBlock(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*output.CidrBlockAssociation.AssociationId, *output.CidrBlockAssociation.CidrBlock, output.CidrBlockAssociation.CidrBlockState.State)
}

func DisassociateVpcCidrBlock() {
	input := ec2.DisassociateVpcCidrBlockInput{
		AssociationId: aws.String("vpc-cidr-assoc-0f5582b9236bb592a"),
	}
	output, err := svc.DisassociateVpcCidrBlock(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.CidrBlockAssociation.CidrBlockState.State)
}

func DeleteVpc() {
	input := ec2.DeleteVpcInput{
		VpcId: aws.String("vpc-006ba5fb389c1ccba"),
	}
	output, err := svc.DeleteVpc(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func CreateSubnet() {
	input := ec2.CreateSubnetInput{
		VpcId:     aws.String("vpc-006ba5fb389c1ccba"),
		CidrBlock: aws.String("10.0.84.0/22"),
	}
	output, err := svc.CreateSubnet(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	subnet := output.Subnet
	displaySubnet(*subnet)
}

func DescribeSubnets() {
	input := ec2.DescribeSubnetsInput{
		Filters: []types.Filter{{
			Name:   aws.String("vpc-id"),
			Values: []string{"vpc-006ba5fb389c1ccba"},
		},
		},
	}
	onput, err := svc.DescribeSubnets(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, subnet := range onput.Subnets {
		displaySubnet(subnet)
		fmt.Println()
	}
}

func CreateSubnetCidrReservation() {
	input := ec2.CreateSubnetCidrReservationInput{
		Cidr:            aws.String("10.0.84.0/23"),
		SubnetId:        aws.String("subnet-0ac563926878c2e84"),
		ReservationType: types.SubnetCidrReservationTypeExplicit,
	}
	output, err := svc.CreateSubnetCidrReservation(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*output.SubnetCidrReservation.SubnetCidrReservationId)
}

func displaySubnet(subnet types.Subnet) {
	fmt.Println(*subnet.AssignIpv6AddressOnCreation,
		*subnet.AvailabilityZone,
		*subnet.AvailabilityZoneId,
		*subnet.AvailableIpAddressCount,
		*subnet.CidrBlock,
		*subnet.DefaultForAz,
		*subnet.EnableDns64,
		*subnet.Ipv6Native,
		*subnet.MapPublicIpOnLaunch,
		*subnet.OwnerId,
		*subnet.PrivateDnsNameOptionsOnLaunch,
		subnet.State,
		*subnet.SubnetArn,
		*subnet.SubnetId,
	)
}
