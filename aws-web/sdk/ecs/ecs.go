package ecs

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/tools"
)

var (
	svc *ecs.Client
)

func init() {
	svc = ecs.NewFromConfig(awsenv.Cfg)
}

func DescribeClusters() {
	input := &ecs.ListClustersInput{}
	output, err := svc.ListClusters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	describeInput := &ecs.DescribeClustersInput{
		Clusters: output.ClusterArns,
	}
	describeOutput, err := svc.DescribeClusters(awsenv.EmptyContext, describeInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, cluster := range describeOutput.Clusters {
		fmt.Println("ClusterName:", *cluster.ClusterName, "\tArn:", *cluster.ClusterArn, "\tStatus:", *cluster.Status)
	}
}

func DescribeCapacityProviders() {
	input := &ecs.DescribeCapacityProvidersInput{}
	output, err := svc.DescribeCapacityProviders(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, provider := range output.CapacityProviders {
		fmt.Println("Name:", *provider.Name, "\tCapacityProviderArn:", *provider.CapacityProviderArn)
	}
}

func CreateAndDeleteCluster() {
	createInput := &ecs.CreateClusterInput{
		CapacityProviders: []string{"FARGATE", "FARGATE_SPOT"},
		ClusterName:       aws.String(awsenv.ECSClusterName),
		// Configuration: &types.ClusterConfiguration{
		// 	ExecuteCommandConfiguration: &types.ExecuteCommandConfiguration{
		// 		Logging: types.ExecuteCommandLoggingOverride,
		// 		LogConfiguration: &types.ExecuteCommandLogConfiguration{
		// 			CloudWatchEncryptionEnabled: true,
		// 			CloudWatchLogGroupName:      aws.String(awsenv.ECSClusterName),
		// 		},
		// 	},
		// },
		DefaultCapacityProviderStrategy: []types.CapacityProviderStrategyItem{{
			CapacityProvider: aws.String("FARGATE"),
		}},
		ServiceConnectDefaults: &types.ClusterServiceConnectDefaultsRequest{
			Namespace: aws.String(awsenv.NamespaceName),
		},
		Settings: []types.ClusterSetting{{
			Name:  types.ClusterSettingNameContainerInsights,
			Value: aws.String("enabled"),
		}},
		Tags: []types.Tag{{
			Key:   aws.String("CreatedBy"),
			Value: aws.String("amazon"),
		}},
	}

	createOutput, err := svc.CreateCluster(awsenv.EmptyContext, createInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("ClusterName:", *createOutput.Cluster.ClusterName, "\tArn:", *createOutput.Cluster.ClusterArn, "\tStatus:", *createOutput.Cluster.Status)

	deleteInput := &ecs.DeleteClusterInput{
		Cluster: createOutput.Cluster.ClusterArn,
	}
	deleteOutput, err := svc.DeleteCluster(awsenv.EmptyContext, deleteInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Status:", *deleteOutput.Cluster.Status)
}

func DescribeTaskDefinition() {
	input := &ecs.ListTaskDefinitionFamiliesInput{
		Status: types.TaskDefinitionFamilyStatusActive,
	}
	output, err := svc.ListTaskDefinitionFamilies(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, task := range output.Families {
		input := &ecs.DescribeTaskDefinitionInput{
			TaskDefinition: &task,
		}
		output, err := svc.DescribeTaskDefinition(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		registerTaskInput := &ecs.RegisterTaskDefinitionInput{
			ContainerDefinitions:    output.TaskDefinition.ContainerDefinitions,
			Family:                  output.TaskDefinition.Family,
			Cpu:                     output.TaskDefinition.Cpu,
			EphemeralStorage:        output.TaskDefinition.EphemeralStorage,
			ExecutionRoleArn:        output.TaskDefinition.ExecutionRoleArn,
			InferenceAccelerators:   output.TaskDefinition.InferenceAccelerators,
			IpcMode:                 output.TaskDefinition.IpcMode,
			Memory:                  output.TaskDefinition.Memory,
			NetworkMode:             output.TaskDefinition.NetworkMode,
			PidMode:                 output.TaskDefinition.PidMode,
			PlacementConstraints:    output.TaskDefinition.PlacementConstraints,
			ProxyConfiguration:      output.TaskDefinition.ProxyConfiguration,
			RequiresCompatibilities: output.TaskDefinition.RequiresCompatibilities,
			RuntimePlatform:         output.TaskDefinition.RuntimePlatform,
			//Tags:                    output.Tags,
			TaskRoleArn: output.TaskDefinition.TaskRoleArn,
			Volumes:     output.TaskDefinition.Volumes,
		}

		registerTaskOutput, err := svc.RegisterTaskDefinition(awsenv.EmptyContext, registerTaskInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Family:", *registerTaskOutput.TaskDefinition.Family, "\tRevision:", registerTaskOutput.TaskDefinition.Revision)
		taskDefinition := *registerTaskOutput.TaskDefinition.Family + ":" + tools.String(registerTaskOutput.TaskDefinition.Revision)
		updateServiceInput := &ecs.UpdateServiceInput{
			Cluster:        aws.String(awsenv.ECSClusterName),
			Service:        &task,
			TaskDefinition: &taskDefinition,
		}
		updateServiceOutput, err := svc.UpdateService(awsenv.EmptyContext, updateServiceInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("\tCurrent:", *updateServiceOutput.Service.TaskDefinition)
	}
}

func ListContainerInstances() {
	input := &ecs.ListContainerInstancesInput{
		Cluster: aws.String(awsenv.ECSClusterName),
	}
	output, err := svc.ListContainerInstances(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ContainerInstanceArns)
}

func ListAccountSettings() {
	input := &ecs.ListAccountSettingsInput{}
	output, err := svc.ListAccountSettings(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, setting := range output.Settings {
		fmt.Println(setting.Name.Values(), *setting.PrincipalArn, *setting.Value)
	}
}

func ListServices() {
	input := &ecs.ListServicesInput{
		Cluster:    aws.String(awsenv.ECSClusterName),
		MaxResults: aws.Int32(100),
	}
	output, err := svc.ListServices(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range output.ServiceArns {
		fmt.Println(v)
	}
}

func ListTaskDefinitions() {
	input := &ecs.ListTaskDefinitionsInput{
		MaxResults:   aws.Int32(100),
		FamilyPrefix: aws.String(awsenv.ECSServiceName),
		Sort:         types.SortOrderDesc,
	}
	output, err := svc.ListTaskDefinitions(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range output.TaskDefinitionArns {
		fmt.Println(v)
	}
}

func ListTasks() {
	input := &ecs.ListTasksInput{
		Cluster:    aws.String(awsenv.ECSClusterName),
		MaxResults: aws.Int32(100),
	}
	output, err := svc.ListTasks(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range output.TaskArns {
		fmt.Println(v)
	}
}

func DescribeServices() {
	input := &ecs.DescribeServicesInput{
		Cluster:  aws.String(awsenv.ECSClusterName),
		Services: []string{awsenv.ECSServiceName},
	}
	fmt.Println(awsenv.ECSServiceName)
	output, err := svc.DescribeServices(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range output.Services {
		fmt.Println(*v.ClusterArn, " ", *v.ServiceName, " ", *v.RoleArn, " ", *v.ServiceArn, " ")
		for _, vv := range v.Deployments {
			fmt.Print("\t", vv.UpdatedAt, "\n")
		}
		for _, vv := range v.LoadBalancers {
			fmt.Println(*vv.TargetGroupArn)
		}

		fmt.Println(v.NetworkConfiguration.AwsvpcConfiguration.SecurityGroups, v.NetworkConfiguration.AwsvpcConfiguration.Subnets)

		for _, vv := range v.ServiceRegistries {
			fmt.Println(*vv.RegistryArn)
		}
		fmt.Println(*v.TaskDefinition)
		fmt.Println()
	}
}

func DescribeTaskDefinition2() {
	input := &ecs.DescribeTaskDefinitionInput{
		TaskDefinition: aws.String(awsenv.ECSServiceName),
	}
	output, err := svc.DescribeTaskDefinition(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	data, _ := json.Marshal(output.TaskDefinition)
	fmt.Println(string(data))
}
