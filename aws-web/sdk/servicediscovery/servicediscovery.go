package servicediscovery

import (
	"fmt"
	"time"

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

	count := 0
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
			// // GetInstance
			// input := &servicediscovery.GetInstanceInput{
			// 	InstanceId: instance.Id,
			// 	ServiceId:  service.Id,
			// }

			// output, err := svc.GetInstance(awsenv.EmptyContext, input)
			// if err != nil {
			// 	fmt.Println(err)
			// 	continue
			// }

			ip := instance.Attributes["AWS_INSTANCE_IPV4"]
			ipint := tools.IP4ToLong(ip) + 1
			instance.Attributes["AWS_INSTANCE_IPV4"] = tools.LongToIP4(int64(ipint))

			registerInput := &servicediscovery.RegisterInstanceInput{
				ServiceId:  input.ServiceId,
				InstanceId: instance.Id,
				Attributes: instance.Attributes,
			}

			registerOutput, err := svc.RegisterInstance(awsenv.EmptyContext, registerInput)

			if err != nil {
				fmt.Println(err)
				continue
			}

			fmt.Println("RegisterInstanceOperationId:", *registerOutput.OperationId)
		}

		if count == 1 {
			continue
		}

		for _, instance := range output.Instances {
			deregisterInput := &servicediscovery.DeregisterInstanceInput{
				ServiceId:  service.Id,
				InstanceId: instance.Id,
			}

			deregisterOutput, err := svc.DeregisterInstance(awsenv.EmptyContext, deregisterInput)
			if err != nil {
				// DuplicateRequest: Another operation of type RegisterInstance
				fmt.Println(err)
				continue
			}

			fmt.Println("DeregisterInstanceOperationId:", deregisterOutput.OperationId)
		}

		count++
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
		Name:        aws.String(awsenv.CloudMapServiceName),
		NamespaceId: aws.String(awsenv.NamespaceId),
	}

	result, err := svc.CreateService(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("arn:", *result.Service.Arn, "\tName:", *result.Service.Name)

	deleteInput := &servicediscovery.DeleteServiceInput{
		Id: result.Service.Id,
	}

	_, err = svc.DeleteService(awsenv.EmptyContext, deleteInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("ok")
}

func GetNamespace() {
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String(awsenv.NamespaceId),
	}

	result, err := svc.GetNamespace(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("arn:", *result.Namespace.Arn,
		"\tName:", *result.Namespace.Name,
		"\tServiceCount:", result.Namespace.ServiceCount,
		"\tHostedZoneId:", *result.Namespace.Properties.DnsProperties.HostedZoneId,
	)
}

func GetNamespaceDetail(namespaceId string) *types.Namespace {
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String(namespaceId),
	}

	result, err := svc.GetNamespace(awsenv.EmptyContext, input)
	if err != nil {
		return nil
	}

	return result.Namespace
}

func CreateNamespace() {
	input := &servicediscovery.CreateHttpNamespaceInput{
		Name: aws.String(awsenv.NamespaceName),
	}

	output, err := svc.CreateHttpNamespace(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("CreateOperationId:", *output.OperationId)

	operationInput := &servicediscovery.GetOperationInput{
		OperationId: output.OperationId,
	}

	namespaceid := ""
	for {
		time.Sleep(time.Second * 2)
		operationOutput, err := svc.GetOperation(awsenv.EmptyContext, operationInput)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Status:", operationOutput.Operation.Status)

		if operationOutput.Operation.Status == types.OperationStatusSuccess {
			namespaceid = operationOutput.Operation.Targets["NAMESPACE"]
			break
		}
	}

	if namespaceid == "" {
		return
	}

	deleteInput := &servicediscovery.DeleteNamespaceInput{
		Id: aws.String(namespaceid),
	}

	deleteOutput, err := svc.DeleteNamespace(awsenv.EmptyContext, deleteInput)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("DeleteOperationId:", *deleteOutput.OperationId)
}
