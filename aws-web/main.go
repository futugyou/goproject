package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/core"
	"github.com/futugyousuzu/goproject/awsgolang/services"
	"github.com/futugyousuzu/goproject/awsgolang/viewmodel"

	// "github.com/futugyousuzu/goproject/awsgolang/sdk/config"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/servicediscovery"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/route53"
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
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/secretsmanager"
	// "github.com/futugyousuzu/goproject/awsgolang/sdk/loadbalancing"
	"github.com/futugyousuzu/goproject/awsgolang/sdk/iot"
)

func main() {
	defer awsenv.DeleteAll()

	if os.Getenv("GITHUB_ACTIONS") == "true" {
		SyncData()
		return
	}

	fmt.Println("not from github action")

	awsServiceDemo()

	awsSdkDemo()

}

func SyncData() {
	parameterService := services.NewParameterService()
	parameterService.SyncAllParameter()

	secService := services.NewEcsClusterService()
	secService.SyncAllEcsServices()

	// current data is useful, no need to update
	// config := services.NewAwsConfigService()
	// config.SyncResourcesByConfig()
}

func awsServiceDemo() {
	// regionService := services.NewRegionService()
	// regions, _ := regionService.GetRegions()
	// fmt.Println(regions)

	// accountService := services.NewAccountService()
	// paging := core.Paging{Page: 1, Limit: 2}
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

	// account := services.UserAccount{
	// 	Id:              "64ca42443b1359c86cb3f144",
	// 	Alias:           "demo10",
	// 	SecretAccessKey: "SecretAccessKey11",
	// 	AccessKeyId:     "AccessKeyId112",
	// 	Region:          "Region",
	// }

	// err := accountService.UpdateAccount(account)
	// fmt.Println(err)

	// err := accountService.DeleteAccount("64ca4766c2d58ba6236c9400")
	// fmt.Println(err)

	// parameterService := services.NewParameterService()
	// parameterService.SyncAllParameter()

	// secService := services.NewEcsClusterService()
	// secService.SyncAllEcsServices()

	// keyService := services.NewKeyValueService()
	// keyService.CreateKeyValue("key", "value")
	// fmt.Println(keyService.GetValueByKey("key"))
	// fmt.Println(keyService.GetAllKeyValues())

	// S3Test()
}

func S3Test() {
	s3Service := services.NewS3bucketService()
	// s3Service.InitData()
	paging := core.Paging{Page: 1, Limit: 5}
	filter := viewmodel.S3BucketFilter{BucketName: ""}
	s3s := s3Service.GetS3Buckets(paging, filter)
	for _, s := range s3s {
		fmt.Println(s.Name)
	}

	f := viewmodel.S3BucketItemFilter{}
	s3items := s3Service.GetS3BucketItems(f)
	for _, s := range s3items {
		fmt.Println(s.Key)
	}
}

func awsSdkDemo() {
	// servicediscovery.ListNamespace()
	// servicediscovery.ListServices()
	// servicediscovery.RegisterInstance()
	// servicediscovery.CreateService()
	// servicediscovery.CreateNamespace()
	// servicediscovery.GetNamespace()

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
	// ecs.ListServices()
	// ecs.ListTaskDefinitions()
	// ecs.ListTasks()
	// ecs.DescribeServices()
	// ecs.DescribeTaskDefinition2()
	// ecs.DescribeTasks()

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

	// names := s3.ListBuckets()
	// for _, name := range names {
	// s3.ListObjectsV2(name)
	// s3.GetBucketCors(name)
	// s3.GetBucketPolicy(name)
	// s3.GetBucketPolicyStatus(name)
	// s3.GetBucketAccelerateConfiguration(name)
	// s3.GetBucketAcl(name)
	// s3.GetBucketLocation(name)
	// s3.GetBucketWebsite(name)
	// }

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

	// config.DeliverConfigSnapshot()
	// config.DescribeConfigRules()
	// config.DescribeConfigurationRecorders()
	// config.DeleteConfigurationRecorder()

	// route53.GetHostedZone()

	// secretsmanager.GetSecretValue()

	// loadbalancing.DescribeTargetGroups()
	Iot()
}

func Iot() {
	// iot.ListJobs()
	// iot.ListThings()
	// iot.ListThingTypes()
	// iot.ListThingGroups()
	// iot.ListThingRegistrationTasks()
	// iot.ListTopicRuleDestinations()
	// iot.ListTopicRules()
	// iot.ListActiveViolations()
	// iot.ListStreams()
	// iot.DescribeThing()
	// iot.DescribeThingGroup()
	// iot.GetRegistrationCode()
	// iot.DescribeEndpoint()
	// iot.ListBillingGroups()
	// iot.ListAuthorizers()
	// iot.ListCACertificates()
	iot.ListCertificates()
}
