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
		VpcId: aws.String("vpc-0664da5448a5dc3f0"),
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
		VpcId:     aws.String("vpc-0664da5448a5dc3f0"),
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

func GetSubnetCidrReservations() {
	input := ec2.GetSubnetCidrReservationsInput{
		SubnetId: aws.String("subnet-0ac563926878c2e84"),
	}
	output, err := svc.GetSubnetCidrReservations(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	ipv4s := output.SubnetIpv4CidrReservations
	for _, r := range ipv4s {
		fmt.Println(*r.Cidr, *r.OwnerId, r.ReservationType, *r.SubnetCidrReservationId, *r.SubnetId)
	}
}

func AssociateSubnetCidrBlock() {
	input := ec2.AssociateSubnetCidrBlockInput{
		Ipv6CidrBlock: aws.String("2001:db8:1234:1a00::/64"),
		SubnetId:      aws.String("subnet-0ac563926878c2e84"),
	}
	output, err := svc.AssociateSubnetCidrBlock(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*output.Ipv6CidrBlockAssociation.AssociationId, *output.Ipv6CidrBlockAssociation.Ipv6CidrBlock, output.Ipv6CidrBlockAssociation.Ipv6CidrBlockState)
}

func DisassociateSubnetCidrBlock() {
	input := ec2.DisassociateSubnetCidrBlockInput{
		AssociationId: aws.String("associate-subnet-0ac563926878c2e84"),
	}
	output, err := svc.DisassociateSubnetCidrBlock(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DeleteSubnetCidrReservation() {
	input := ec2.DeleteSubnetCidrReservationInput{
		SubnetCidrReservationId: aws.String("scr-0eee42ae462d56a1a"),
	}
	output, err := svc.DeleteSubnetCidrReservation(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DeleteSubnet() {
	input := ec2.DeleteSubnetInput{
		SubnetId: aws.String("subnet-0ac563926878c2e84"),
	}
	output, err := svc.DeleteSubnet(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DescribeNetworkAcls() {
	input := ec2.DescribeNetworkAclsInput{}
	output, err := svc.DescribeNetworkAcls(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, acl := range output.NetworkAcls {
		displayAcl(acl)
	}
}

func CreateNetworkAcl() {
	input := ec2.CreateNetworkAclInput{
		VpcId: aws.String("vpc-0664da5448a5dc3f0"),
	}
	output, err := svc.CreateNetworkAcl(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	displayAcl(*output.NetworkAcl)
}

func CreateNetworkAclEntry() {
	input := ec2.CreateNetworkAclEntryInput{
		NetworkAclId: aws.String("acl-00419da83c857ade9"),
		Egress:       aws.Bool(true),
		RuleAction:   types.RuleActionAllow,
		RuleNumber:   aws.Int32(100),
		Protocol:     aws.String("88"),
		PortRange: &types.PortRange{
			From: aws.Int32(77),
			To:   aws.Int32(77),
		},
		CidrBlock: aws.String("10.0.0.0/16"),
	}
	output, err := svc.CreateNetworkAclEntry(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DeleteNetworkAclEntry() {
	input := ec2.DeleteNetworkAclEntryInput{
		RuleNumber:   aws.Int32(100),
		NetworkAclId: aws.String("acl-00419da83c857ade9"),
		Egress:       aws.Bool(true),
	}
	output, err := svc.DeleteNetworkAclEntry(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DeleteNetworkAcl() {
	input := ec2.DeleteNetworkAclInput{
		NetworkAclId: aws.String("acl-00419da83c857ade9"),
	}
	output, err := svc.DeleteNetworkAcl(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DescribeNatGateways() {
	input := ec2.DescribeNatGatewaysInput{}
	output, err := svc.DescribeNatGateways(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, nat := range output.NatGateways {
		displayNatgateway(nat)
	}
}

func CreateNatGateway() {
	input := ec2.CreateNatGatewayInput{
		SubnetId:         aws.String("subnet-0505218e5192b90be"),
		ConnectivityType: types.ConnectivityTypePrivate,
	}
	output, err := svc.CreateNatGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	displayNatgateway(*output.NatGateway)
}

func DeleteNatGateway() {
	input := ec2.DeleteNatGatewayInput{
		NatGatewayId: aws.String("nat-0f7863130dd048954"),
	}
	output, err := svc.DeleteNatGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*output.NatGatewayId)
}

func DescribeInternetGateways() {
	input := ec2.DescribeInternetGatewaysInput{}
	output, err := svc.DescribeInternetGateways(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, gateway := range output.InternetGateways {
		displayInternetGateway(gateway)
	}
}

func CreateInternetGateway() {
	input := ec2.CreateInternetGatewayInput{}
	output, err := svc.CreateInternetGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	displayInternetGateway(*output.InternetGateway)
}

func AttachInternetGateway() {
	input := ec2.AttachInternetGatewayInput{
		InternetGatewayId: aws.String("igw-078cf9209f27c2fc5"),
		VpcId:             aws.String("vpc-0664da5448a5dc3f0"),
	}
	output, err := svc.AttachInternetGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DetachInternetGateway() {
	input := ec2.DetachInternetGatewayInput{
		InternetGatewayId: aws.String("igw-078cf9209f27c2fc5"),
		VpcId:             aws.String("vpc-0664da5448a5dc3f0"),
	}
	output, err := svc.DetachInternetGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}

func DeleteInternetGateway() {
	input := ec2.DeleteInternetGatewayInput{
		InternetGatewayId: aws.String("igw-078cf9209f27c2fc5"),
	}
	output, err := svc.DeleteInternetGateway(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
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

func displayAcl(acl types.NetworkAcl) {
	fmt.Println(*acl.VpcId, *acl.OwnerId, *acl.NetworkAclId, *acl.IsDefault)
	for _, a := range acl.Associations {
		fmt.Println("\t", *a.SubnetId, *a.NetworkAclId, *a.NetworkAclAssociationId)
	}
	for _, v := range acl.Entries {
		fmt.Println("\t", *v.CidrBlock, *v.Egress, *v.Protocol, v.RuleAction, *v.RuleNumber)
	}
}

func displayNatgateway(nat types.NatGateway) {
	fmt.Print(
		*nat.NatGatewayId, "\t",
		nat.ConnectivityType, "\t",
		nat.CreateTime, "\t",
		nat.State,
	)
	for _, v := range nat.NatGatewayAddresses {
		if v.AllocationId != nil {
			fmt.Print("\t", *v.AllocationId)
		}
		if v.NetworkInterfaceId != nil {
			fmt.Print("\t", *v.NetworkInterfaceId)
		}
		if v.PrivateIp != nil {
			fmt.Print("\t", *v.PrivateIp)
		}
		if v.PublicIp != nil {
			fmt.Print("\t", *v.PublicIp)
		}
	}
	fmt.Println()
}

func displayInternetGateway(gateway types.InternetGateway) {
	fmt.Print(*gateway.InternetGatewayId, "\t", *gateway.OwnerId, "\t")
	if len(gateway.Attachments) == 0 {
		fmt.Println()
	}
	for _, v := range gateway.Attachments {
		fmt.Println(*v.VpcId, v.State)
	}
}
