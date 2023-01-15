package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/futugyousuzu/goproject/awsgolang/awsenv"
)

var (
	svc *ssm.Client
)

func init() {
	svc = ssm.NewFromConfig(awsenv.Cfg)
}

func ListAssociations() {
	input := &ssm.ListAssociationsInput{}
	output, err := svc.ListAssociations(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, association := range output.Associations {
		fmt.Println(&association.AssociationName)
	}
}

func ListCommands() {
	input := &ssm.ListCommandsInput{}
	output, err := svc.ListCommands(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, command := range output.Commands {
		fmt.Println(&command.Comment)
	}
}

func GetParametersByPath() {
	input := &ssm.GetParametersByPathInput{
		Path: aws.String("/"),
	}
	output, err := svc.GetParametersByPath(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, parameter := range output.Parameters {
		fmt.Println(*parameter.Name)
	}
}

func GetParameters(name string) {
	input := &ssm.GetParametersInput{
		Names: []string{name},
	}
	output, err := svc.GetParameters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range output.Parameters {
		fmt.Println("Value:\n", *p.Value)
	}
}

func DescribeParameters() {
	input := &ssm.DescribeParametersInput{
		// max value 50
		MaxResults: aws.Int32(50),
	}
	output, err := svc.DescribeParameters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, p := range output.Parameters {
		fmt.Println("Name:", *p.Name, "\tDataType:", *p.DataType, "\tTier:", p.Tier, "\tType:", p.Type, "\tPolicies:", p.Policies)
		GetParameters(*p.Name)
		fmt.Println()
	}
}
