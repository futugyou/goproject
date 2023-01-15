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
