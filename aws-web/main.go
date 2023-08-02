package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/services"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/servicediscovery"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/cloudwatch"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/cloudwatchlogs"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/ecr"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/iam"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/ecs"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/ssm"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/ec2"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/s3"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/dynamodb"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/efs"
)

func main() {
	defer awsenv.DeleteAll()

	// regionService := services.NewRegionService()
	// regions, _ := regionService.GetAllRegionInCurrentAccount()
	// fmt.Println(regions)

	accountService := services.NewAccountService()
	// accountService.AccountInit()
	// paging := core.Paging{Page: 10, Limit: 2}
	// accounts := accountService.GetAccountsByPaging(paging)
	// fmt.Println(accounts)

	// account := services.UserAccount{
	// 	Alias:           "demo10",
	// 	SecretAccessKey: "SecretAccessKey",
	// 	AccessKeyId:     "AccessKeyId",
	// 	Region:          "Region",
	// }

	// err := accountService.CreateAccount(account)
	// fmt.Println(err)

	account := services.UserAccount{
		Id:              "64ca42443b1359c86cb3f144",
		Alias:           "demo10",
		SecretAccessKey: "SecretAccessKey11",
		AccessKeyId:     "AccessKeyId112",
		Region:          "Region",
	}

	err := accountService.UpdateAccount(account)
	fmt.Println(err)
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
	// iam.CreateLoginProfile()
	// iam.UpdateLoginProfile()
	// iam.AddUserToGroup()

	// ecs.DescribeClusters()
	// ecs.DescribeCapacityProviders()
	// ecs.CreateAndDeleteCluster()
	// ecs.DescribeTaskDefinition()
	// ecs.ListContainerInstances()
	// ecs.ListAccountSettings()

	// ssm.DescribeParameters()
	// ssm.ListAssociations()
	// ssm.ListCommands()
	// ssm.GetParametersByPath()
	// ssm.GetParameters("/")
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
	// ec2.DescribeNatGateways()
	// ec2.CreateNatGateway()
	// ec2.DeleteNatGateway()
	// ec2.DescribeInternetGateways()
	// ec2.CreateInternetGateway()
	// ec2.AttachInternetGateway()
	// ec2.DetachInternetGateway()
	// ec2.DeleteInternetGateway()

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

	// efs.DescribeFileSystems()
	// efs.DescribeAccessPoints()
}
