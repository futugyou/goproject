package servicediscoverydemo

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/servicediscovery"
)

func ListServices() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := servicediscovery.New(sess)
	input := &servicediscovery.ListServicesInput{}

	result, err := svc.ListServices(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func CreateService() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := servicediscovery.New(sess)
	input := &servicediscovery.CreateServiceInput{
		DnsConfig: &servicediscovery.DnsConfig{
			DnsRecords: []*servicediscovery.DnsRecord{
				{
					TTL:  aws.Int64(60),
					Type: aws.String("A"),
				},
			},
			RoutingPolicy: aws.String("WEIGHTED"),
		},
		Name:        aws.String("oneapp-feedback"),
		NamespaceId: aws.String("ns-3n5se3ecgsvtqebi"),
	}

	result, err := svc.CreateService(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			case servicediscovery.ErrCodeResourceLimitExceeded:
				fmt.Println(servicediscovery.ErrCodeResourceLimitExceeded, aerr.Error())
			case servicediscovery.ErrCodeNamespaceNotFound:
				fmt.Println(servicediscovery.ErrCodeNamespaceNotFound, aerr.Error())
			case servicediscovery.ErrCodeServiceAlreadyExists:
				fmt.Println(servicediscovery.ErrCodeServiceAlreadyExists, aerr.Error())
			case servicediscovery.ErrCodeTooManyTagsException:
				fmt.Println(servicediscovery.ErrCodeTooManyTagsException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func GetNamespace() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := servicediscovery.New(sess)
	input := &servicediscovery.GetNamespaceInput{
		Id: aws.String("ns-3n5se3ecgsvtqebi"),
	}

	result, err := svc.GetNamespace(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			case servicediscovery.ErrCodeNamespaceNotFound:
				fmt.Println(servicediscovery.ErrCodeNamespaceNotFound, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func ListNamespace() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := servicediscovery.New(sess)
	input := &servicediscovery.ListNamespacesInput{}

	result, err := svc.ListNamespaces(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}

func CreateNamespace() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(endpoints.ApSoutheast1RegionID)},
	}))
	svc := servicediscovery.New(sess)
	input := &servicediscovery.CreateHttpNamespaceInput{
		Name: aws.String("gateway.com"),
	}

	result, err := svc.CreateHttpNamespace(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case servicediscovery.ErrCodeInvalidInput:
				fmt.Println(servicediscovery.ErrCodeInvalidInput, aerr.Error())
			case servicediscovery.ErrCodeNamespaceAlreadyExists:
				fmt.Println(servicediscovery.ErrCodeNamespaceAlreadyExists, aerr.Error())
			case servicediscovery.ErrCodeResourceLimitExceeded:
				fmt.Println(servicediscovery.ErrCodeResourceLimitExceeded, aerr.Error())
			case servicediscovery.ErrCodeDuplicateRequest:
				fmt.Println(servicediscovery.ErrCodeDuplicateRequest, aerr.Error())
			case servicediscovery.ErrCodeTooManyTagsException:
				fmt.Println(servicediscovery.ErrCodeTooManyTagsException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
