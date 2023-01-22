package main

import (
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"

	// "github.com/futugyousuzu/goproject/awsgolang/servicediscovery"
	// "github.com/futugyousuzu/goproject/awsgolang/cloudwatch"
	// "github.com/futugyousuzu/goproject/awsgolang/cloudwatchlogs"
	// "github.com/futugyousuzu/goproject/awsgolang/ecr"
	// "github.com/futugyousuzu/goproject/awsgolang/iam"
	// "github.com/futugyousuzu/goproject/awsgolang/ecs"
	// "github.com/futugyousuzu/goproject/awsgolang/ssm"
	"github.com/futugyousuzu/goproject/awsgolang/ec2"
	// "github.com/futugyousuzu/goproject/awsgolang/s3"
	// "github.com/futugyousuzu/goproject/awsgolang/dynamodb"
)

func main() {
	defer awsenv.DeleteAll()

	// servicediscovery.ListNamespace()
	// servicediscovery.ListServices()
	// servicediscovery.RegisterInstance()
	// servicediscovery.CreateService()
	// servicediscovery.CreateNamespace()

	// cloudwatch.GetMetricData()
	// cloudwatch.ListDashboards()
	// cloudwatch.ListMetrics()
	// cloudwatch.GetMetricStatistics()

	// cloudwatchlogs.DescribeExportTasks()
	// cloudwatchlogs.DescribeLogGroups()
	// cloudwatchlogs.GetLogEvents()
	// cloudwatchlogs.DescribeLogStreams()
	// cloudwatchlogs.GetLogGroupFields()
	// cloudwatchlogs.DescribeQueries()

	// ecr.DescribeRepositories()
	// ecr.CreateRepository()
	// ecr.DeleteRepository()

	// iam.ListUsers()
	// iam.ListAccessKeys("")
	// iam.ListGroups()
	// iam.ListAccountAliases()
	// iam.CreateAccountAlias()
	// iam.DeleteAccountAlias()
	// iam.ListInstanceProfiles()
	// iam.ListPolicies()
	// iam.ListRoles()
	// iam.CreateGroup()
	// iam.DeleteGroup()
	// iam.CreateUser()
	// iam.DeleteUser()
	// iam.GetAccountAuthorizationDetails()

	// ecs.DescribeClusters()
	// ecs.DescribeCapacityProviders()
	// ecs.CreateAndDeleteCluster()
	// ecs.DescribeTaskDefinition()
	// ecs.ListContainerInstances()
	// ecs.ListAccountSettings()

	// ssm.ListAssociations()
	// ssm.ListCommands()
	// ssm.GetParametersByPath()
	// ssm.GetParameters("/")
	// ssm.DescribeParameters()
	// ssm.PutParameter()
	// ssm.DeleteParameter()

	// ec2.DescribeSecurityGroups()
	// ec2.DescribeVpcs()
	// ec2.CreateVpc()
	// ec2.AssociateVpcCidrBlock()
	// ec2.DisassociateVpcCidrBlock()
	// ec2.DeleteVpc()
	// ec2.CreateSubnet()
	// ec2.DescribeSubnets()
	// ec2.CreateSubnetCidrReservation()
	// ec2.GetSubnetCidrReservations()
	// ec2.AssociateSubnetCidrBlock()
	// ec2.DisassociateSubnetCidrBlock()
	// ec2.DeleteSubnetCidrReservation()
	// ec2.DeleteSubnet()
	// ec2.DescribeNetworkAcls()
	// ec2.CreateNetworkAcl()
	// ec2.CreateNetworkAclEntry()
	// ec2.DeleteNetworkAclEntry()
	// ec2.DeleteNetworkAcl()
	ec2.DescribeNatGateways()

	// s3.ListBuckets()
	// s3.ListObjectsV2("/")
	// s3.GetObject("", "")
	// s3.PutObject("", "")

	// dynamodb.ListGlobalTables()
	// dynamodb.ListTables()
	// dynamodb.DescribeTable("/")
	// dynamodb.CreateTable()
	// dynamodb.CreateBackup()
	// dynamodb.ListBackups()
	// dynamodb.DeleteBackup("/")
	// dynamodb.DeleteTable()
	// dynamodb.Scan("/")
	// dynamodb.Query("/")
	// dynamodb.GetItem()
	// dynamodb.PutItem()
	// dynamodb.UpdateItem()
	// dynamodb.DeleteItem()
}
