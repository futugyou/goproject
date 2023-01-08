package servicediscovery

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery"
	"github.com/aws/aws-sdk-go-v2/service/servicediscovery/types"

	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
	"github.com/futugyousuzu/goproject/awsgolang/tools"
)

var (
	svc *servicediscovery.Client
)

func init() {
	svc = servicediscovery.NewFromConfig(awsenv.Cfg)
}

func ListNamespace() {
	input := &servicediscovery.ListNamespacesInput{}

	result, err := svc.ListNamespaces(awsenv.EmptyContext, input)
	if err != nil {

		fmt.Println(err.Error())
		return
	}

	for _, namespace := range result.Namespaces {
		fmt.Println("arn:", *namespace.Arn, "\tId", *namespace.Id, "\tName:", *namespace.Name)
	}
}

func ListServices() {
	input := &servicediscovery.ListServicesInput{}

	result, err := svc.ListServices(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, service := range result.Services {
		fmt.Print("Name:", *service.Name, "\tId:", *service.Id)

		// ListInstances
		input := &servicediscovery.ListInstancesInput{
			ServiceId: service.Id,
		}
		output, err := svc.ListInstances(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, instance := range output.Instances {
			fmt.Println("\t\tInstanceID:", *instance.Id, "\tIP:", instance.Attributes["AWS_INSTANCE_IPV4"])
			// for key, value := range instance.Attributes {
			// 	// instance.Attributes["AWS_INSTANCE_IPV4"]
			// 	fmt.Println("\tkey:", key, "\tvalue:", value)
			// }

			// GetInstance
			// input := &servicediscovery.GetInstanceInput{
			// 	InstanceId: instance.Id,
			// 	ServiceId:  service.Id,
			// }

			// output, err := svc.GetInstance(awsenv.EmptyContext, input)
			// if err != nil {
			// 	fmt.Println(err)
			// 	continue
			// }

			// for key, value := range output.Instance.Attributes {
			// 	fmt.Println("\tkey:", key, "\tvalue:", value)
			// }
		}
	}
}

func RegisterInstance() {
	input := &servicediscovery.ListServicesInput{}

	result, err := svc.ListServices(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, service := range result.Services {
		input := &servicediscovery.ListInstancesInput{
			ServiceId: service.Id,
		}
		output, err := svc.ListInstances(awsenv.EmptyContext, input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, instance := range output.Instances {
			// GetInstance
			input := &servicediscovery.GetInstanceInput{
				InstanceId: instance.Id,
				ServiceId:  service.Id,
			}

			output, err := svc.GetInstance(awsenv.EmptyContext, input)
			if err != nil {
				fmt.Println(err)
				continue
			}

			ip := output.Instance.Attributes["AWS_INSTANCE_IPV4"]
			ipint := tools.IP4ToLong(ip) + 1
			output.Instance.Attributes["AWS_INSTANCE_IPV4"] = tools.LongToIP4(int64(ipint))

			registerInput := &servicediscovery.RegisterInstanceInput{
				ServiceId:  input.ServiceId,
				InstanceId: input.InstanceId,
				Attributes: output.Instance.Attributes,
			}

			registerOutput, err := svc.RegisterInstance(awsenv.EmptyContext, registerInput)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("OperationId:", *registerOutput.OperationId)
		}
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
