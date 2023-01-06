package servicediscovery

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *servicediscovery.Client
)

func init() {
	svc = servicediscovery.NewFromConfig(awsenv.Cfg)
}

func ListServices() {
	input := &servicediscovery.ListServicesInput{}

	result, err := svc.ListServices(awsenv.EmptyContext, input)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, service := range result.Services {
		fmt.Println("arn:", *service.Arn, "\tName:", *service.Name)
	}
}

func CreateService() {
	input := &servicediscovery.CreateServiceInput{
		DnsConfig: &types.DnsConfig{
			DnsRecords: []types.DnsRecord{
				{
					TTL:  aws.Int64(60),
					Type: types.RecordTypeA,
				},
			},
			RoutingPolicy: types.RoutingPolicyWeighted,
		},
		Name:        aws.String("oneapp-feedback"),
		NamespaceId: aws.String("ns-3n5se3ecgsvtqebi"),
	}

	result, err := svc.CreateService(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("arn:", *result.Service.Arn, "\tName:", *result.Service.Name)
}

func GetNamespace() {
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String("ns-3n5se3ecgsvtqebi"),
	}

	result, err := svc.GetNamespace(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("arn:", *result.Namespace.Arn, "\tName:", *result.Namespace.Name, "\tServiceCount:", result.Namespace.ServiceCount)
}

func ListNamespace() {
	input := &servicediscovery.ListNamespacesInput{}

	result, err := svc.ListNamespaces(awsenv.EmptyContext, input)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, namespace := range result.Namespaces {
		fmt.Println("arn:", *namespace.Arn, "\tName:", *namespace.Name, "\tServiceCount:", namespace.ServiceCount)
	}
}

func CreateNamespace() {
	input := &servicediscovery.CreateHttpNamespaceInput{
		Name: aws.String("gateway.com"),
	}

	result, err := svc.CreateHttpNamespace(awsenv.EmptyContext, input)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	fmt.Println(*result.OperationId)
}
