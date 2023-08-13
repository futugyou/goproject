package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
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
		Path:           aws.String("/"),
		WithDecryption: aws.Bool(true),
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
		Names:          []string{name}, // max count 10
		WithDecryption: aws.Bool(true),
	}
	output, err := svc.GetParameters(awsenv.EmptyContext, input)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, p := range output.Parameters {
		fmt.Println("Version:\n", p.Version)
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

func PutParameter() {
	putInput := &ssm.PutParameterInput{
		Name:      aws.String("/Terrform/Configuration/NetworkPolicy"),
		Value:     aws.String("{\"apiVersion\":{\"kind\":\"PodSecurityConfiguration\",\"defaults\":{\"enforce\":\"baseline\",\"enforce-version\":\"latest\",\"audit\":\"restricted\",\"audit-version\":\"latest\",\"warn\":\"restricted\",\"warn-version\":\"latest\"}}}"),
		Overwrite: aws.Bool(true),
		Type:      types.ParameterTypeString,
	}
	putOutput, err := svc.PutParameter(awsenv.EmptyContext, putInput)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Name:", &putInput.Name, "\tTier:", putOutput.Tier, "\tVersion:", putOutput.Version)
	GetParameters(*putInput.Name)
}

func DeleteParameter() {
	input := ssm.DeleteParameterInput{
		Name: aws.String("/Terrform/Configuration/NetworkPolicy"),
	}
	output, err := svc.DeleteParameter(awsenv.EmptyContext, &input)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(output.ResultMetadata)
}
